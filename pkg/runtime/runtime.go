package runtime

import (
	"context"
	"log"

	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/utils"

	_ "github.com/DaoCloud-OpenSource/ferry-proxy/init"
	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/config"
	"github.com/DaoCloud-OpenSource/ferry-proxy/pkg/registries"
)

type Runtime struct {
	ctx           context.Context
	serviceCancel map[string]context.CancelFunc
}

func NewRuntime(ctx context.Context) *Runtime {
	return &Runtime{
		ctx:           ctx,
		serviceCancel: map[string]context.CancelFunc{},
	}
}

func (r *Runtime) OnEndpoint(endpoint *config.EndpointConfig) {
	dialer, err := registries.StreamDialerRegistry.New("dialer", endpoint.Endpoint)
	if err != nil {
		log.Println(err)
		return
	}
	registries.StreamDialerInstance.Put(endpoint.Name, dialer)
}

func (r *Runtime) OnCluster(cluster *config.ClusterConfig) {
	if len(cluster.Cluster) == 0 {
		log.Println("cluster has no endpoints")
		return
	}

	// TODO: support LB
	endpoint := cluster.Cluster[0]
	handler, err := registries.StreamHandlerRegistry.New("forward", endpoint)
	if err != nil {
		log.Println(err)
		return
	}
	registries.StreamHandlerInstance.Put(cluster.Name, handler)
}

func (r *Runtime) OnRouter(router *config.RouterConfig) {
	if len(router.Router) == 0 {
		log.Println("router has no endpoints")
		return
	}
	handler, err := registries.StreamHandlerRegistry.New("http_hosts", router.Router)
	if err != nil {
		log.Println(err)
		return
	}
	registries.StreamHandlerInstance.Put(router.Name, handler)
}

func (r *Runtime) OnListener(listener *config.ListenerConfig) {
	listen, err := registries.StreamListenConfigRegistry.New("listen_config", listener.Listener)
	if err != nil {
		log.Println(err)
		return
	}
	registries.StreamListenConfigInstance.Put(listener.Name, listen)

	if lastCancel, ok := r.serviceCancel[listener.Name]; ok && lastCancel != nil {
		defer lastCancel()
	}

	ctx, cancel := context.WithCancel(r.ctx)
	r.serviceCancel[listener.Name] = cancel

	started := make(chan struct{})
	go r.startService(ctx, listener.Name, listener.Router, started)
	<-started
}

func (r *Runtime) startService(ctx context.Context, listenConfigName string, routerName string, started chan struct{}) {

	listenConfig := registries.StreamListenConfigInstance.Get(listenConfigName)
	if listenConfig == nil {
		log.Printf("listenConfig %q error: not found", listenConfigName)
		return
	}
	router := registries.StreamHandlerInstance.Get(routerName)
	if router == nil {
		log.Printf("handler %q error: not found", routerName)
		return
	}

	listen, err := listenConfig.Listen(r.ctx)
	if err != nil {
		log.Printf("listenConfig %q error: %v", listenConfigName, err)
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		listen.Close()
	}()

	if started != nil {
		close(started)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			if utils.IsClosedConnError(err) {
				log.Printf("listenConfig %q accept connect error: %v", listenConfigName, err)
			}
			return
		}
		go router.ServeStream(r.ctx, conn)
	}
}

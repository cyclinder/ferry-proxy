package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Source struct {
	ClusterCallback  func(cluster *ClusterConfig)
	EndpointCallback func(endpoint *EndpointConfig)
	ListenerCallback func(listener *ListenerConfig)
	RouterCallback   func(router *RouterConfig)
}

func (s *Source) Filter(path string) error {
	if !strings.HasSuffix(path, ".json") {
		return fmt.Errorf("config %q not has suffix .json", path)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load %q: %w", path, err)
	}
	switch {
	case strings.HasSuffix(path, "cluster.json"):
		if s.ClusterCallback == nil {
			return fmt.Errorf("config %q: no cluster callback", path)
		}
		cluster := &ClusterConfig{}
		err := json.Unmarshal(body, cluster)
		if err != nil {
			return err
		}
		s.ClusterCallback(cluster)
	case strings.HasSuffix(path, "endpoint.json"):
		if s.EndpointCallback == nil {
			return fmt.Errorf("config %q: no endpoint callback", path)
		}
		endpoint := &EndpointConfig{}
		err := json.Unmarshal(body, endpoint)
		if err != nil {
			return err
		}
		s.EndpointCallback(endpoint)
	case strings.HasSuffix(path, "listener.json"):
		if s.ListenerCallback == nil {
			return fmt.Errorf("config %q: no listener callback", path)
		}
		listener := &ListenerConfig{}
		err := json.Unmarshal(body, listener)
		if err != nil {
			return err
		}
		s.ListenerCallback(listener)
	case strings.HasSuffix(path, "router.json"):
		if s.RouterCallback == nil {
			return fmt.Errorf("config %q: no router callback", path)
		}
		router := &RouterConfig{}
		err := json.Unmarshal(body, router)
		if err != nil {
			return err
		}
		s.RouterCallback(router)
	}
	return nil
}

func (s *Source) Run(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		log.Println("first load file:", path)
		err = s.Filter(path)
		if err != nil {
			log.Println(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = s.WatchConfig(dir)
	return err
}

func (s *Source) WatchConfig(dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			log.Println("event:", event)
			if event.Op&fsnotify.Write != 0 ||
				event.Op&fsnotify.Create != 0 {
				log.Println("modified file:", event.Name)
				err := s.Filter(event.Name)
				if err != nil {
					log.Println(err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println("error:", err)
		}
	}
}

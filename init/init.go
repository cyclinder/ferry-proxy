package init

import (
	_ "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/filters/stream_dialer/dialer"
	_ "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/filters/stream_handler/forward"
	_ "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/filters/stream_handler/http_direct"
	_ "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/filters/stream_handler/http_hosts"
	_ "github.com/DaoCloud-OpenSource/ferry-proxy/pkg/filters/stream_listen_config/listen_config"
)

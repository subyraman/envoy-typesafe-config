package main

import (
	"fmt"
	"log"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/protobuf/ptypes/any"

	//"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"

	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"

	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"

	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
)

// static_resources:
//   listeners:
//   - name: listener_0
//     address:
//       socket_address:
//         protocol: TCP
//         address: 0.0.0.0
//         port_value: 10000
//     filter_chains:
//     - filters:
//       - name: envoy.filters.network.http_connection_manager
//         typed_config:
//           "@type": type.googleapis.com/envoy.config.filter.network.http_connection_manager.v2.HttpConnectionManager
//           stat_prefix: ingress_http
//           route_config:
//             name: local_route
//             virtual_hosts:
//             - name: local_service
//               domains: ["*"]
//               routes:
//               - match:
//                   prefix: "/"
//                 direct_response:
//                   status: 200
//           http_filters:
//           - name: envoy.filters.http.router

func main() {
	httpRouterFilter := envoy_extensions_filters_network_http_connection_manager_v3.HttpFilter{
		Name: "envoy.filters.http.router",
	}

	simplePathSpecifier := envoy_config_route_v3.RouteMatch_Prefix{
		Prefix: "/",
	}
	simpleRouteMatch := envoy_config_route_v3.RouteMatch{
		PathSpecifier: &simplePathSpecifier,
	}

	simpleDirectResponseAction := envoy_config_route_v3.DirectResponseAction{
		Status: 200,
	}
	simpleRouteAction := envoy_config_route_v3.Route_DirectResponse{
		DirectResponse: &simpleDirectResponseAction,
	}

	simpleRoute := envoy_config_route_v3.Route{
		Match:  &simpleRouteMatch,
		Action: &simpleRouteAction,
	}

	localServiceVirtualHost := envoy_config_route_v3.VirtualHost{
		Name:    "local_service",
		Domains: []string{"*"},
		Routes:  []*envoy_config_route_v3.Route{&simpleRoute},
	}

	routeConfig := envoy_config_route_v3.RouteConfiguration{
		Name:         "local_route",
		VirtualHosts: []*envoy_config_route_v3.VirtualHost{&localServiceVirtualHost},
	}

	httpRouteConfig := envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
		RouteConfig: &routeConfig,
	}

	httpConnectionManager := envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
		StatPrefix:     "ingress_http",
		RouteSpecifier: &httpRouteConfig,
		HttpFilters:    []*envoy_extensions_filters_network_http_connection_manager_v3.HttpFilter{&httpRouterFilter},
	}

	serialized, _ := proto.Marshal(&httpConnectionManager)

	httpConnectionManagerAny := &any.Any{
		TypeUrl: "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
		Value:   serialized,
	}

	// create filter
	httpConnectionManagerTypedConfig := envoy_config_listener_v3.Filter_TypedConfig{
		TypedConfig: httpConnectionManagerAny,
	}

	httpConnectionManagerFilter := envoy_config_listener_v3.Filter{
		Name:       "envoy.filters.network.http_connection_manager",
		ConfigType: &httpConnectionManagerTypedConfig,
	}

	filters := []*envoy_config_listener_v3.Filter{&httpConnectionManagerFilter}

	filterChain := envoy_config_listener_v3.FilterChain{
		Filters: filters,
	}

	socketPortSpecifier := envoy_config_core_v3.SocketAddress_PortValue{
		PortValue: 10000,
	}

	socketAddress := envoy_config_core_v3.SocketAddress{
		Protocol:      envoy_config_core_v3.SocketAddress_TCP,
		Address:       "0.0.0.0",
		PortSpecifier: &socketPortSpecifier,
	}

	addressSocketAddress := envoy_config_core_v3.Address_SocketAddress{
		SocketAddress: &socketAddress,
	}
	address := envoy_config_core_v3.Address{
		Address: &addressSocketAddress,
	}

	listener := envoy_config_listener_v3.Listener{
		Name:         "listener1",
		FilterChains: []*envoy_config_listener_v3.FilterChain{&filterChain},
		Address:      &address,
	}

	staticResources := envoy_config_bootstrap_v3.Bootstrap_StaticResources{
		Listeners: []*envoy_config_listener_v3.Listener{&listener},
	}

	bootstrap := envoy_config_bootstrap_v3.Bootstrap{
		StaticResources: &staticResources,
	}

	foo, err := protojson.Marshal(&bootstrap)
	if err != nil {
		log.Fatalf("error when marshalling: %s", err)
	}

	fmt.Println(string(foo))
}

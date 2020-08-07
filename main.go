package main

import (
	"log"
	"os"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/encoding/protojson"

	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
)

func main() {
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

	httpRouterFilter := envoy_extensions_filters_network_http_connection_manager_v3.HttpFilter{
		Name: "envoy.filters.http.router",
	}

	httpConnectionManager := envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
		StatPrefix:     "ingress_http",
		RouteSpecifier: &httpRouteConfig,
		HttpFilters:    []*envoy_extensions_filters_network_http_connection_manager_v3.HttpFilter{&httpRouterFilter},
	}

	// serialize the `httpConnectionManager` message into a raw byte array
	serialized, _ := proto.Marshal(&httpConnectionManager)

	// create an Any protocol buffer representation
	httpConnectionManagerAny := &any.Any{
		TypeUrl: "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager",
		Value:   serialized,
	}

	// use the Any protobuf message as the typed configuration
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

	socketAddress := envoy_config_core_v3.SocketAddress{
		Protocol: envoy_config_core_v3.SocketAddress_TCP,
		Address:  "0.0.0.0",
		PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
			PortValue: 10000,
		},
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

	err := bootstrap.Validate()
	if err != nil {
		log.Fatalf("Error in validation: %s", err)
	}

	opts := protojson.MarshalOptions{
		Indent: "   ",
	}

	out, err := opts.Marshal(&bootstrap)
	if err != nil {
		log.Fatalf("error when marshalling: %s", err)
	}

	f, err := os.Create("envoy.json")
	if err != nil {
		log.Fatalf("error when creating file: %s", err)
	}

	f.WriteString(string(out))
	f.Close()
}

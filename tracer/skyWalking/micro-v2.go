/**
 * Description:micro-v2版本skyWalking插件
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/3/19
 * Time: 15:03
 * Mail: hammercui@163.com
 *
 */
package skyWalking

import (
	"context"
	"errors"
	"fmt"
	"github.com/hammercui/go2sky"
	"github.com/hammercui/go2sky/propagation"
	language_agent "github.com/hammercui/go2sky/reporter/grpc/language-agent"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"strings"
	"time"
)

const (
	componentIDGoMicroClient = 5008
	componentIDGoMicroServer = 5009
)

var errTracerIsNil = errors.New("tracer is nil")

type clientWrapper struct {
	client.Client

	sw         *go2sky.Tracer
	reportTags []string
}

// ClientOption allow optional configuration of Client
type ClientOption func(*clientWrapper)

// WithClientWrapperReportTags customize span tags
func WithClientWrapperReportTags(reportTags ...string) ClientOption {
	return func(c *clientWrapper) {
		c.reportTags = append(c.reportTags, reportTags...)
	}
}

// Call is used for client calls
func (s *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	span, err := s.sw.CreateExitSpan(ctx, name, req.Service(), func(header string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[propagation.Header] = header
		ctx = metadata.NewContext(ctx, md)
		return nil
	})
	if err != nil {
		return err
	}

	span.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(language_agent.SpanLayer_RPCFramework)

	defer span.End()
	for _, k := range s.reportTags {
		if v, ok := metadata.Get(ctx, k); ok {
			span.Tag(go2sky.Tag(k), v)
		}
	}
	if err = s.Client.Call(ctx, req, rsp, opts...); err != nil {
		span.Error(time.Now(), err.Error())
	}
	return err
}

// Stream is used streaming
func (s *clientWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	span, err := s.sw.CreateExitSpan(ctx, name, req.Service(), func(header string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[propagation.Header] = header
		ctx = metadata.NewContext(ctx, md)
		return nil
	})
	if err != nil {
		return nil, err
	}

	span.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(language_agent.SpanLayer_RPCFramework)

	defer span.End()
	for _, k := range s.reportTags {
		if v, ok := metadata.Get(ctx, k); ok {
			span.Tag(go2sky.Tag(k), v)
		}
	}
	stream, err := s.Client.Stream(ctx, req, opts...)
	if err != nil {
		span.Error(time.Now(), err.Error())
	}
	return stream, err
}

// Publish is used publish message to subscriber
func (s *clientWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	name := fmt.Sprintf("Pub to %s", p.Topic())
	span, err := s.sw.CreateExitSpan(ctx, name, p.ContentType(), func(header string) error {
		mda, _ := metadata.FromContext(ctx)
		md := metadata.Copy(mda)
		md[propagation.Header] = header
		ctx = metadata.NewContext(ctx, md)
		return nil
	})
	if err != nil {
		return err
	}

	span.SetComponent(componentIDGoMicroClient)
	span.SetSpanLayer(language_agent.SpanLayer_RPCFramework)

	defer span.End()
	for _, k := range s.reportTags {
		if v, ok := metadata.Get(ctx, k); ok {
			span.Tag(go2sky.Tag(k), v)
		}
	}
	if err = s.Client.Publish(ctx, p, opts...); err != nil {
		span.Error(time.Now(), err.Error())
	}
	return err
}

// NewClientWrapper accepts a go2sky Tracer and returns a Client Wrapper
func NewClientWrapper(sw *go2sky.Tracer, options ...ClientOption) client.Wrapper {
	return func(c client.Client) client.Client {
		co := &clientWrapper{
			sw:     sw,
			Client: c,
		}
		for _, option := range options {
			option(co)
		}
		return co
	}
}

// NewCallWrapper accepts an go2sky Tracer and returns a Call Wrapper
func NewCallWrapper(sw *go2sky.Tracer, reportTags ...string) client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			if sw == nil {
				return errTracerIsNil
			}

			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			span, err := sw.CreateExitSpan(ctx, name, req.Service(), func(header string) error {
				mda, _ := metadata.FromContext(ctx)
				md := metadata.Copy(mda)
				md[propagation.Header] = header
				ctx = metadata.NewContext(ctx, md)
				return nil
			})
			if err != nil {
				return err
			}

			span.SetComponent(componentIDGoMicroClient)
			span.SetSpanLayer(language_agent.SpanLayer_RPCFramework)

			defer span.End()
			for _, k := range reportTags {
				if v, ok := metadata.Get(ctx, k); ok {
					span.Tag(go2sky.Tag(k), v)
				}
			}
			if err = cf(ctx, node, req, rsp, opts); err != nil {
				span.Error(time.Now(), err.Error())
			}
			return err
		}
	}
}

// NewSubscriberWrapper accepts a go2sky Tracer and returns a Handler Wrapper
func NewSubscriberWrapper(sw *go2sky.Tracer, reportTags ...string) server.SubscriberWrapper {
	return func(next server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			if sw == nil {
				return errTracerIsNil
			}

			name := "Sub from " + msg.Topic()
			span, err := sw.CreateExitSpan(ctx, name, msg.ContentType(), func(header string) error {
				mda, _ := metadata.FromContext(ctx)
				md := metadata.Copy(mda)
				md[propagation.Header] = header
				ctx = metadata.NewContext(ctx, md)
				return nil
			})
			if err != nil {
				return err
			}

			span.SetComponent(componentIDGoMicroClient)
			span.SetSpanLayer(language_agent.SpanLayer_RPCFramework)

			defer span.End()
			for _, k := range reportTags {
				if v, ok := metadata.Get(ctx, k); ok {
					span.Tag(go2sky.Tag(k), v)
				}
			}
			if err = next(ctx, msg); err != nil {
				span.Error(time.Now(), err.Error())
			}
			return err
		}
	}
}

// NewHandlerWrapper accepts a go2sky Tracer and returns a Subscriber Wrapper
func NewHandlerWrapper(sw *go2sky.Tracer, reportTags ...string) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			if sw == nil {
				return errTracerIsNil
			}

			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			span, ctx, err := sw.CreateEntrySpan(ctx, name, func() (string, error) {
				str, _ := metadata.Get(ctx, strings.Title(propagation.Header))
				return str, nil
			})
			if err != nil {
				return err
			}

			span.SetComponent(componentIDGoMicroServer)
			span.SetSpanLayer(language_agent.SpanLayer_RPCFramework)

			defer span.End()
			for _, k := range reportTags {
				if v, ok := metadata.Get(ctx, k); ok {
					span.Tag(go2sky.Tag(k), v)
				}
			}
			if err = fn(ctx, req, rsp); err != nil {
				span.Error(time.Now(), err.Error())
			}
			return err
		}
	}
}
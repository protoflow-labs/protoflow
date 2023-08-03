package grpc

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func SerializeType(m protoreflect.ProtoMessage) (*desc.MessageDescriptor, error) {
	return desc.WrapMessage(m.ProtoReflect().Descriptor())
}

type TypeResolver struct {
	DescLookup map[string]protoreflect.MessageDescriptor
	EnumLookup map[string]protoreflect.EnumDescriptor
}

type SerializedTypeResolver struct {
	DescLookup map[string]*descriptorpb.DescriptorProto
	EnumLookup map[string]*descriptorpb.EnumDescriptorProto
}

// TODO breadchris maybe a better name is type serializer?
func NewTypeResolver() *TypeResolver {
	return &TypeResolver{
		DescLookup: map[string]protoreflect.MessageDescriptor{},
		EnumLookup: map[string]protoreflect.EnumDescriptor{},
	}
}

func NewSerializedTypeResolver() *SerializedTypeResolver {
	return &SerializedTypeResolver{
		DescLookup: map[string]*descriptorpb.DescriptorProto{},
		EnumLookup: map[string]*descriptorpb.EnumDescriptorProto{},
	}
}

func (t *TypeResolver) ResolveLookup(m protoreflect.ProtoMessage) *TypeResolver {
	nt := *t
	dl, el := ResolveTypeLookup(m.ProtoReflect().Descriptor(), nt.DescLookup, nt.EnumLookup)
	for k, v := range dl {
		nt.DescLookup[k] = v
	}
	for k, v := range el {
		nt.EnumLookup[k] = v
	}
	return &nt
}

func (t *TypeResolver) Serialize() *SerializedTypeResolver {
	sr := NewSerializedTypeResolver()
	for k, v := range t.DescLookup {
		m, err := desc.WrapMessage(v)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to wrap message %s", k)
			continue
		}
		sr.DescLookup[k] = m.AsDescriptorProto()
	}
	for k, v := range t.EnumLookup {
		e, err := desc.WrapEnum(v)
		if err != nil {
			log.Warn().Err(err).Msgf("unable to wrap enum %s", k)
			continue
		}
		sr.EnumLookup[k] = e.AsEnumDescriptorProto()
	}
	return sr
}

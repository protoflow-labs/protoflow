package grpc

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
	"testing"
)

// TODO breadchris figure out if we can base a graph off of just protobuf files
func TestParseProto(t *testing.T) {
	res, err := ParseProto("")
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]interface{}{}
	for _, r := range res {
		buildUninterpretedMapForFile(r.GetFile(), m)
		for _, s := range r.GetServices() {
			t.Log(s.GetFullyQualifiedName())
			for _, m := range s.GetMethods() {
				t.Log(m.GetFullyQualifiedName())
				mo := m.GetMethodOptions()
				mo.GetUninterpretedOption()
				t.Log(mo.String())
				for _, o := range mo.GetUninterpretedOption() {
					for _, n := range o.GetName() {
						t.Log(n)
					}
				}
			}
		}
	}
}

func buildUninterpretedMapForFile(fd *desc.FileDescriptor, opts map[string]interface{}) {
	buildUninterpretedMap(fd.GetName(), fd.GetFileOptions().GetUninterpretedOption(), opts)
	for _, md := range fd.GetMessageTypes() {
		buildUninterpretedMapForMessage(fd.GetPackage(), md.AsDescriptorProto(), opts)
	}
	//for _, extd := range fd.GetExtensions() {
	//	buildUninterpretedMap(qualify(fd.GetPackage(), extd.GetName()), extd.GetOp.GetUninterpretedOption(), opts)
	//}
	for _, ed := range fd.GetEnumTypes() {
		buildUninterpretedMapForEnum(fd.GetPackage(), ed.AsEnumDescriptorProto(), opts)
	}
	for _, sd := range fd.GetServices() {
		buildUninterpretedMap(sd.GetFullyQualifiedName(), sd.GetServiceOptions().GetUninterpretedOption(), opts)
		for _, mtd := range sd.GetMethods() {
			buildUninterpretedMap(mtd.GetFullyQualifiedName(), mtd.GetMethodOptions().GetUninterpretedOption(), opts)
		}
	}
}

func buildUninterpretedMapForMessage(qual string, md *descriptorpb.DescriptorProto, opts map[string]interface{}) {
	fqn := qualify(qual, md.GetName())
	buildUninterpretedMap(fqn, md.GetOptions().GetUninterpretedOption(), opts)
	for _, fld := range md.GetField() {
		buildUninterpretedMap(qualify(fqn, fld.GetName()), fld.GetOptions().GetUninterpretedOption(), opts)
	}
	for _, ood := range md.GetOneofDecl() {
		buildUninterpretedMap(qualify(fqn, ood.GetName()), ood.GetOptions().GetUninterpretedOption(), opts)
	}
	for _, extr := range md.GetExtensionRange() {
		buildUninterpretedMap(qualify(fqn, fmt.Sprintf("%d-%d", extr.GetStart(), extr.GetEnd()-1)), extr.GetOptions().GetUninterpretedOption(), opts)
	}
	for _, nmd := range md.GetNestedType() {
		buildUninterpretedMapForMessage(fqn, nmd, opts)
	}
	for _, extd := range md.GetExtension() {
		buildUninterpretedMap(qualify(fqn, extd.GetName()), extd.GetOptions().GetUninterpretedOption(), opts)
	}
	for _, ed := range md.GetEnumType() {
		buildUninterpretedMapForEnum(fqn, ed, opts)
	}
}

func buildUninterpretedMapForEnum(qual string, ed *descriptorpb.EnumDescriptorProto, opts map[string]interface{}) {
	fqn := qualify(qual, ed.GetName())
	buildUninterpretedMap(fqn, ed.GetOptions().GetUninterpretedOption(), opts)
	for _, evd := range ed.GetValue() {
		buildUninterpretedMap(qualify(fqn, evd.GetName()), evd.GetOptions().GetUninterpretedOption(), opts)
	}
}

type ident string
type aggregate string

func buildUninterpretedMap(prefix string, uos []*descriptorpb.UninterpretedOption, opts map[string]interface{}) {
	for _, uo := range uos {
		parts := make([]string, len(uo.GetName()))
		for i, np := range uo.GetName() {
			if np.GetIsExtension() {
				parts[i] = fmt.Sprintf("(%s)", np.GetNamePart())
			} else {
				parts[i] = np.GetNamePart()
			}
		}
		uoName := fmt.Sprintf("%s:%s", prefix, strings.Join(parts, "."))
		key := uoName
		i := 0
		for {
			if _, ok := opts[key]; !ok {
				break
			}
			i++
			key = fmt.Sprintf("%s#%d", uoName, i)
		}
		var val interface{}
		switch {
		case uo.AggregateValue != nil:
			val = aggregate(uo.GetAggregateValue())
		case uo.IdentifierValue != nil:
			val = ident(uo.GetIdentifierValue())
		case uo.DoubleValue != nil:
			val = uo.GetDoubleValue()
		case uo.PositiveIntValue != nil:
			val = int(uo.GetPositiveIntValue())
		case uo.NegativeIntValue != nil:
			val = int(uo.GetNegativeIntValue())
		default:
			val = string(uo.GetStringValue())
		}
		opts[key] = val
	}
}

func qualify(base, name string) string {
	if base == "" {
		return name
	}
	return base + "." + name
}

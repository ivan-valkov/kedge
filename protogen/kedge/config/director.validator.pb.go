// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kedge/config/director.proto

package kedge_config

import fmt "fmt"
import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/golang/protobuf/proto"
import math "math"
import _ "github.com/mwitkow/go-proto-validators"
import  _ "github.com/improbable-eng/kedge/protogen/kedge/config/grpc/routes"
import  _ "github.com/improbable-eng/kedge/protogen/kedge/config/http/routes"
import  _ "github.com/improbable-eng/kedge/protogen/kedge/config/http/routes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *DirectorConfig) Validate() error {
	if nil == this.Grpc {
		return github_com_mwitkow_go_proto_validators.FieldError("Grpc", fmt.Errorf("message must exist"))
	}
	if this.Grpc != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Grpc); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Grpc", err)
		}
	}
	if nil == this.Http {
		return github_com_mwitkow_go_proto_validators.FieldError("Http", fmt.Errorf("message must exist"))
	}
	if this.Http != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Http); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Http", err)
		}
	}
	return nil
}
func (this *DirectorConfig_Grpc) Validate() error {
	for _, item := range this.Routes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Routes", err)
			}
		}
	}
	return nil
}
func (this *DirectorConfig_Http) Validate() error {
	for _, item := range this.Routes {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Routes", err)
			}
		}
	}
	for _, item := range this.AdhocRules {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("AdhocRules", err)
			}
		}
	}
	return nil
}
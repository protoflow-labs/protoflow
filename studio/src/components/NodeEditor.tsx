import {InputEditor} from "@/components/blockEditors/InputEditor";
import {GenericNodeEditor} from "@/components/blockEditors/GenericNodeEditor";
import React from "react";
import {Node as ProtoNode} from "@/rpc/graph_pb";
import {Bucket, Collection, Config, Function, GRPC, Prompt, Query, REST, Template, File, Route} from "@/rpc/block_pb";
import {Message} from "@bufbuild/protobuf";

type NodeEditorProps = {
    node: ProtoNode | null;
};

export function NodeEditor(props: NodeEditorProps) {
    if (!props.node) {
        return null;
    }

    switch (props.node.config.case) {
        case "input":
            return <InputEditor node={props.node} />;
        case "collection":
            return <GenericNodeEditor node={props.node} nodeConfig={'collection'} nodeConfigType={Collection} />
        case "query":
            return <GenericNodeEditor node={props.node} nodeConfig={'query'} nodeConfigType={Query} />
        case "function":
            return <GenericNodeEditor node={props.node} nodeConfig={"function"} nodeConfigType={Function} />;
        case "bucket":
            return <GenericNodeEditor node={props.node} nodeConfig={"bucket"} nodeConfigType={Bucket} />;
        case "rest":
            return <GenericNodeEditor node={props.node} nodeConfig={"rest"} nodeConfigType={REST} />;
        case "grpc":
            return <GenericNodeEditor node={props.node} nodeConfig={"grpc"} nodeConfigType={GRPC} />;
        case "prompt":
            return <GenericNodeEditor node={props.node} nodeConfig={"prompt"} nodeConfigType={Prompt} />;
        case "configuration":
            return <GenericNodeEditor node={props.node} nodeConfig={"configuration"} nodeConfigType={Config} />;
        case "template":
            return <GenericNodeEditor node={props.node} nodeConfig={"template"} nodeConfigType={Template} />;
        case "file":
            return <GenericNodeEditor node={props.node} nodeConfig={"file"} nodeConfigType={File} />;
        case "route":
            return <GenericNodeEditor node={props.node} nodeConfig={"route"} nodeConfigType={Route} />;
        default:
            return null;
    }
}

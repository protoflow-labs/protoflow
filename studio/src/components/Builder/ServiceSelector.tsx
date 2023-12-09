import { useProjectContext } from "@/providers/ProjectProvider";
import React, {useCallback, useEffect, useState} from "react";
import {GRPCMethod, GRPCService, NodeExecution, ProviderState } from "@/rpc/project_pb";
import { Node } from "@/rpc/graph_pb";
import {Button, Divider, Dropdown, Input, Option} from "@fluentui/react-components";
import {useEditorContext} from "@/providers/EditorProvider";
import {projectService} from "@/lib/api";
import {List, MessageBar, Stack} from "@fluentui/react";
import {toast} from "react-hot-toast";
import {SiOpenstack} from "react-icons/si";
import {MethodInputForm} from "@/components/Builder/MethodInputForm";
import {JsonViewer} from "@/components/jsonViewer";

export function ServiceSelector() {
    const [service, setService] = useState<GRPCService|undefined>(undefined);
    const [method, setMethod] = useState<GRPCMethod|undefined>(undefined);
    const [output, setOutput] = useState<NodeExecution[]>([]);

    const [methodName, setMethodName] = useState<string>('');

    const [host, setHost] = useState<string>('http://localhost:8080');
    const [services, setServices] = useState<GRPCService[]>([]);

    const getTypeInfo = () => {
        projectService.getGRPCServerInfo({
            host: host,
        }).then((res) => {
            setServices(res.services ?? []);
        }).catch((err) => {
            toast.error(err.message);
        })
    }

    useEffect(() => {
        getTypeInfo();
    }, []);

    const addMethod = () => {
        if (!service) {
            return;
        }
        projectService.addMethod({
            file: service?.file,
            service: service?.name,
            method: methodName,
            package: service?.package,
        }).then((res) => {
            toast.success('Method added to ' + service?.name);
        }).catch((err) => {
            toast.error(err.message);
        })
    }

    const runMethod = async (data: any) => {
        try {
            setOutput([]);
            const stream = projectService.runGRPCMethod({
                package: service?.package,
                service: service?.name,
                method: method?.name,
                host: host,
                input: JSON.stringify(data),
            })
            for await (const exec of stream) {
                console.log(exec);
                setOutput((prev) => {
                    return [...prev, exec];
                });
            }
        } catch (e: any) {
            toast.error(e.toString(), {
                duration: 10000,
            });
        }
    }

    return (
        <Stack>
            <Stack.Item>
                <Input value={host} placeholder={"Host"} onChange={(e) => {
                    setHost(e.target.value);
                }} />
                <Button onClick={getTypeInfo}>Connect</Button>
            </Stack.Item>
            <Stack.Item>
                <Dropdown
                    aria-labelledby={"service-selector"}
                    placeholder={"Select a service"}
                    value={service?.name || ''}
                    onOptionSelect={(ev, data) => {
                        const idx = parseInt(data.optionValue ?? "0");
                        setService(services[idx]);
                    }}
                >
                    {services.map((ep, idx) => {
                        return (
                            <Option key={idx} value={idx.toString()}>{ep.name}</Option>
                        )
                    })}
                </Dropdown>
                {service && (
                    <Dropdown
                        aria-labelledby={"method-selector"}
                        placeholder={"Select a method"}
                        value={method?.name || ''}
                        onOptionSelect={(ev, data) => {
                            const idx = parseInt(data.optionValue ?? "0");
                            setMethod(service.methods[idx]);
                        }}
                    >
                        {
                            service.methods.map((m, idx) => {
                                return (
                                    <Option key={idx} value={idx.toString()}>{m.name}</Option>
                                )
                            })
                        }
                    </Dropdown>
                )}
            </Stack.Item>
            <Stack.Item>
                {method && (
                    <>
                        {method.typeInfo ? (
                            <MethodInputForm typeInfo={method.typeInfo} onSubmit={(data: any) => {
                                void runMethod(data);
                            }} />
                        ) : (
                            <p>No type info</p>
                        )}
                    </>
                )}
            </Stack.Item>
            <Stack.Item>
                <List items={output.map(o => o.output)} onRenderCell={(item?: string) => {
                    const getText = () => {
                        if (!item) {
                            return '';
                        }
                        const parsed = JSON.parse(item);
                        const p = JSON.parse(parsed.output);
                        // if (Object.keys(p).length === 1) {
                        //     const value: any = Object.values(p)[0];
                        //     return (<ReactMarkdown>{value.toString()}</ReactMarkdown>);
                        // }
                        return <JsonViewer data={p} />;
                    }
                    return (
                        <MessageBar messageBarType={0}>{getText()}</MessageBar>
                    );
                }} />
            </Stack.Item>
            {service && (
                <Stack.Item>
                    <Input value={methodName} onChange={(e) => {
                        setMethodName(e.target.value);
                    }} />
                    <Button onClick={addMethod}>Add Method</Button>
                </Stack.Item>
            )}
        </Stack>
    );
}

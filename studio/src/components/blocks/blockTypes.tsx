import React, {ReactNode} from "react";
import { TbBucket } from "react-icons/tb";
import {HiCircleStack, HiOutlineMagnifyingGlass, HiPencilSquare} from "react-icons/hi2";
import { HiCodeBracket } from "react-icons/hi2";
import {MdHttp, MdLightbulb, MdOutbound} from "react-icons/md";
import {AiOutlineMail} from "react-icons/ai";
import {BiClipboard} from "react-icons/bi";


// The most basic metadata about a type of block, used by both the sidebar and the main canvas
export interface BlockType {
    label: string;
    typeName: string;
    image: ReactNode;
}

// export interface BlockConfigs {
//     bucket: Record<never, unknown>;
//     collection: { entity?: { collection: string } };
// }


// TODO breadchris now that project types exist, this should be able to be generically defined
export const blockTypes: BlockType[] = [
    {
        label: "Bucket",
        typeName:'bucket',
        image: <TbBucket className="h-5 w-5 bg-gray-800" />,
    },
    {
        label: "Collection",
        typeName:'collection',
        image: <HiCircleStack className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "Function",
        typeName:'function',
        image: <HiCodeBracket className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "GRPC",
        typeName:'grpc',
        image:<MdOutbound className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "Manual Input",
        typeName:'input',
        image: <HiPencilSquare className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "DB Query",
        typeName:'query',
        image: <HiOutlineMagnifyingGlass className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "Queue",
        typeName:'queue',
        image: <AiOutlineMail className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "REST",
        typeName:'rest',
        image: <MdHttp className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "Prompt",
        typeName: 'prompt',
        image: <MdLightbulb className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "Config",
        typeName: 'configuration',
        image: <BiClipboard className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "Template",
        typeName: 'template',
        image: <BiClipboard className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "File",
        typeName: 'file',
        image: <BiClipboard className="h-5 w-5 bg-gray-800" />
    },
    {
        label: "Route",
        typeName: 'route',
        image: <BiClipboard className="h-5 w-5 bg-gray-800" />
    }
    ]

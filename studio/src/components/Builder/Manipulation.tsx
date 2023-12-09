import * as React from "react";
import {
    FlatTree,
    FlatTreeItem,
    TreeItemLayout,
    TreeOpenChangeData,
    TreeOpenChangeEvent,
    HeadlessFlatTreeItemProps,
    useHeadlessFlatTree_unstable,
    TreeItemValue,
    FlatTreeItemProps,
} from "@fluentui/react-components";
import { Delete20Regular, MoreHorizontal20Regular, Edit20Regular } from "@fluentui/react-icons";
import {
    Button,
    Menu,
    MenuItem,
    MenuList,
    MenuPopover,
    MenuTrigger,
    useRestoreFocusTarget,
} from "@fluentui/react-components";
import {useProjectContext} from "@/providers/ProjectProvider";
import { ProviderState } from "@/rpc/project_pb";
import {useEditorContext} from "@/providers/EditorProvider";

type ItemProps = HeadlessFlatTreeItemProps & { content: string };

const subtrees: ItemProps[][] = [
    [
        { value: "1", content: "Level 1, item 1" },
        { value: "1-1", parentValue: "1", content: "Item 1-1" },
        { value: "1-2", parentValue: "1", content: "Item 1-2" },
    ],

    [
        { value: "2", content: "Level 1, item 2" },
        { value: "2-1", parentValue: "2", content: "Item 2-1" },
    ],
];

type CustomTreeItemProps = FlatTreeItemProps & {
    onRemoveItem?: (value: string) => void;
};

const ActionsExample: React.FC<{isItemRemovable: boolean, handleRemoveItem: () => void}> = ({ isItemRemovable, handleRemoveItem}) => {
    return (
        <>
            {isItemRemovable && (
                <Button
                    aria-label="Remove item"
                    appearance="subtle"
                    onClick={handleRemoveItem}
                    icon={<Delete20Regular />}
                />
            )}
            <Button aria-label="Edit" appearance="subtle" icon={<Edit20Regular />} />
            <Menu>
                <MenuTrigger disableButtonEnhancement>
                    <Button
                        aria-label="More options"
                        appearance="subtle"
                        icon={<MoreHorizontal20Regular />}
                    />
                </MenuTrigger>

                <MenuPopover>
                    <MenuList>
                        <MenuItem>New</MenuItem>
                        <MenuItem>New Window</MenuItem>
                        <MenuItem disabled>Open File</MenuItem>
                        <MenuItem>Open Folder</MenuItem>
                    </MenuList>
                </MenuPopover>
            </Menu>
        </>
    )
};

const CustomTreeItem = React.forwardRef(
    (
        { onRemoveItem, ...props }: CustomTreeItemProps,
        ref: React.Ref<HTMLDivElement> | undefined
    ) => {
        const { nodeLookup } = useProjectContext();
        const { setSelectedNodes } = useEditorContext();
        const focusTargetAttribute = useRestoreFocusTarget();
        const level = props["aria-level"];
        const value = props.value as string;
        const isItemRemovable = level !== 1 && !value.endsWith("-btn");

        const handleRemoveItem = React.useCallback(() => {
            onRemoveItem?.(value);
        }, [value, onRemoveItem]);

        const selectNode = () => {
            const node = nodeLookup[value];
            if (node) {
                setSelectedNodes([node]);
            }
        }

        return (
            <Menu positioning="below-end" openOnContext>
                <MenuTrigger disableButtonEnhancement>
                    <FlatTreeItem
                        aria-description="has actions"
                        {...focusTargetAttribute}
                        {...props}
                        ref={ref}
                        onClick={selectNode}
                    >
                        <TreeItemLayout
                            actions={<ActionsExample isItemRemovable={isItemRemovable} handleRemoveItem={handleRemoveItem} /> }
                        >
                            {props.children}
                        </TreeItemLayout>
                    </FlatTreeItem>
                </MenuTrigger>
                <MenuPopover>
                    <MenuList>
                        <MenuItem onClick={handleRemoveItem}>Remove item</MenuItem>
                    </MenuList>
                </MenuPopover>
            </Menu>
        );
    }
);

export const Manipulation = () => {
    const { providers } = useProjectContext();

    const providerSubtree: ItemProps[][] = providers.map((r, index) => {
        const p = r.provider;
        if (!p) {
            return [];
        }

        const parentContent: ItemProps = {
            value: p.id,
            content: p.name,
        }

        const resError = r.info && r.info.state === ProviderState.ERROR;
        return [parentContent, ...r.nodes.map(node => ({
            value: node.id,
            parentValue: p.id,
            content: node.name,
        }))];
    });

    const [trees, setTrees] = React.useState(providerSubtree);
    const itemToFocusRef = React.useRef<HTMLDivElement>(null);
    const [itemToFocusValue, setItemToFocusValue] =
        React.useState<TreeItemValue>();

    const handleOpenChange = (
        event: TreeOpenChangeEvent,
        data: TreeOpenChangeData
    ) => {
        // casting here to string as no number values are used in this example
        const value = data.value as string;
        if (value.endsWith("-btn")) {
            const subtreeIndex = Number(value[0]) - 1;
            addFlatTreeItem(subtreeIndex);
        }
    };

    const addFlatTreeItem = (subtreeIndex: number) =>
        setTrees((currentTrees) => {
            const lastItem =
                currentTrees[subtreeIndex][currentTrees[subtreeIndex].length - 1];
            const newItemValue = `${subtreeIndex + 1}-${
                Number(lastItem.value.toString().slice(2)) + 1
            }`;
            setItemToFocusValue(newItemValue);
            const nextSubTree: ItemProps[] = [
                ...currentTrees[subtreeIndex],
                {
                    value: newItemValue,
                    parentValue: currentTrees[subtreeIndex][0].value,
                    content: `New item ${newItemValue}`,
                },
            ];

            return [
                ...currentTrees.slice(0, subtreeIndex),
                nextSubTree,
                ...currentTrees.slice(subtreeIndex + 1),
            ];
        });

    const removeFlatTreeItem = React.useCallback(
        (value: string) =>
            setTrees((currentTrees) => {
                const subtreeIndex = Number(value[0]) - 1;
                const currentSubTree = trees[subtreeIndex];
                const itemIndex = currentSubTree.findIndex(
                    (item) => item.value === value
                );
                const nextSubTree = trees[subtreeIndex].filter(
                    (_item, index) => index !== itemIndex
                );

                const nextItemValue = currentSubTree[itemIndex + 1]?.value;
                const prevItemValue = currentSubTree[itemIndex - 1]?.value;
                setItemToFocusValue(nextItemValue || prevItemValue);

                return [
                    ...currentTrees.slice(0, subtreeIndex),
                    nextSubTree,
                    ...currentTrees.slice(subtreeIndex + 1),
                ];
            }),
        [trees]
    );

    const flatTree = useHeadlessFlatTree_unstable(
        React.useMemo(
            () => {
                return providerSubtree.length > 0 ? providerSubtree.reduce((acc, t) => {
                    return [...acc, ...t];
                }, [] as ItemProps[]): []
            },

            [providers, providerSubtree]
        ),
        { defaultOpenItems: [], onOpenChange: handleOpenChange }
    );

    React.useEffect(() => {
        if (itemToFocusRef.current) {
            itemToFocusRef.current.focus();
            setItemToFocusValue(undefined);
        }
    }, [itemToFocusValue]);

    return (
        <FlatTree {...flatTree.getTreeProps()} aria-label="Manipulation">
            {Array.from(flatTree.items(), (item) => {
                const { content, ...treeItemProps } = item.getTreeItemProps();
                return (
                    <CustomTreeItem
                        {...treeItemProps}
                        key={item.value}
                        onRemoveItem={removeFlatTreeItem}
                        ref={item.value === itemToFocusValue ? itemToFocusRef : undefined}
                    >
                        {content}
                    </CustomTreeItem>
                );
            })}
        </FlatTree>
    );
};
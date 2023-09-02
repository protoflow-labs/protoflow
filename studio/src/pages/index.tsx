import { useProjectContext } from "@/providers/ProjectProvider";

import "reactflow/dist/style.css";
import {ChatPage} from "@/components/Chat/Chat";

export default function Home() {
  const { project } = useProjectContext();

  if (!project) {
    return null;
  }

  return (
    <main className="flex">
      <ChatPage />
    </main>
  );
}

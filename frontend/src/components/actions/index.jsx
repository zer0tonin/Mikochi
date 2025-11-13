import { h } from "preact";
import Download from "./download";
import Rename from "./rename";
import Delete from "./delete";
import CopyLink from "./copylink";

export const DirActions = ({filePath, refresh, setRefresh}) => {
  if (window.innerWidth < 768) {
    console.log("here");
  }

  return (
    <>
      <Download filePath={`${filePath}/`} />
      <Rename
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
      <Delete
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
    </>
  );
}

export const FileActions = ({filePath, refresh, setRefresh}) => {
  return (
    <>
      <Download filePath={filePath} />
      <CopyLink filePath={filePath} />
      <Rename
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
      <Delete
        filePath={filePath}
        refresh={refresh}
        setRefresh={setRefresh}
      />
    </>
  );
}

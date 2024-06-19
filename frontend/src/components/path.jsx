import { h } from "preact";
import { useLocation } from "preact-iso";
import Icon from "./icon";

export const Path = ({ fileInfo, currentDir }) => {
  const location = useLocation();

  // a Path can represent multiple sub-directories during search
  const splitPath = fileInfo.path.split("/");

  if (splitPath.length == 1) {
    // browse
    if (fileInfo.isDir) {
      const target =
        currentDir == "/"
          ? `/${fileInfo.path}/`
          : `${currentDir}/${fileInfo.path}/`;
      return (
        <a href="#" onClick={() => location.route(target)}>
          {fileInfo.path}
        </a>
      );
    }
    return <span>{fileInfo.path}</span>;
  }

  // search
  return (
    <>
      {splitPath
        .map((val, i) => {
          if (i == splitPath.length - 1 && !fileInfo.isDir) {
            return <span key={i}>{val}</span>;
          }
          const target = `${splitPath.slice(0, i + 1).join("/")}/`;
          return (
            <a href="#" key={i} onClick={() => location.route(target)}>
              {val}
            </a>
          );
        })
        .reduce(
          (acc, val) => (acc === null ? [val] : [...acc, " / ", val]),
          null,
        )}
    </>
  );
};

export const FolderIcon = ({ currentDir, fileInfo }) => {
  const location = useLocation();

  const target =
    currentDir == "/"
      ? `${fileInfo.path.startsWith("/") ? "" : "/"}${fileInfo.path}/`
      : `${currentDir}/${fileInfo.path}/`;
  console.log(target);

  return <Icon name="folder" onClick={() => location.route(target)} />;
};

export const DoubleDotPath = ({ currentDir }) => {
  const location = useLocation();
  const split = currentDir.split("/");
  const target = split.slice(0, split.length - 1).join("/");

  return (
    <a href="#" onClick={() => location.route(target == "" ? "/" : target)}>
      ..
    </a>
  );
};

export const DoubleDotFolderIcon = ({ currentDir }) => {
  const location = useLocation();

  const split = currentDir.split("/");
  const target = split.slice(0, split.length - 1).join("/");

  return (
    <Icon
      name="folder"
      onClick={() => location.route(target == "" ? "/" : target)}
    />
  );
};

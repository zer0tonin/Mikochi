import { h } from "preact";
import { useLocation } from "preact-iso";
import Icon from "./icon";

export const Path = ({ fileInfo, currentDir, isSearch }) => {
  const location = useLocation();


  if (!isSearch) {
    const fileName = currentDir === "/" ? fileInfo.path.slice(currentDir.length) : fileInfo.path.slice(currentDir.length + 1);
    // browse
    if (fileInfo.isDir) {
      return <a onClick={() => location.route(fileInfo.path)}>{fileName}</a>;
    }
    return <span>{fileName}</span>;
  }

  // a Path can represent multiple sub-directories during search
  const splitPath = fileInfo.path.split("/");

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
            <a key={i} onClick={() => location.route(target)}>
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

export const FolderIcon = ({ fileInfo }) => {
  const location = useLocation();

  return <Icon name="folder" onClick={() => location.route(fileInfo.path)} />;
};

export const DoubleDotPath = ({ currentDir }) => {
  const location = useLocation();
  const split = currentDir.split("/");
  const target = split.slice(0, split.length - 1).join("/");

  return <a onClick={() => location.route(target == "" ? "/" : target)}>..</a>;
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

import { h } from "preact";

export const Path = ({ fileInfo, currentDir }) => {
  // a Path can represent multiple sub-directories during search
  const splitPath = fileInfo.path.split("/");

  return (
    <>
      {splitPath
        .map((val, i) => {
          if (i == splitPath.length - 1 && !fileInfo.isDir) {
            // files are shown as text and not links
            return <span key={i}>{val}</span>;
          }
          // preact-router has trouble with relative links so we rebuild it from the root
          const target = currentDir != "" ? `/${currentDir}/${splitPath.slice(0, i + 1).join("/")}/` : `/${splitPath.slice(0, i + 1).join("/")}/`;
          return (
            <a
              key={i}
              href={target}
            >
              {val}
            </a>
          );
        })
        .reduce(
          (acc, val) => (acc === null ? [val] : [...acc, " / ", val]),
          null
        )}
    </>
  );
};

export const DoubleDotPath = ({ currentDir }) => {
  const split = currentDir.split("/")
  const target = split.slice(0, split.length - 1).join("/")
  // same as above, link as to start with / for preact-router
  return (
    <a
      href={`/${target}`}
    >
      ..
    </a>
  );
}

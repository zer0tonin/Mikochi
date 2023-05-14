import { h } from "preact";
import { route } from "preact-router";

const Path = ({ fileInfo }) => {
  const splitPath = fileInfo.path.split("/");

  return (
    <>
      {splitPath
        .map((val, i) => {
          if (i == splitPath.length - 1 && !fileInfo.isDir) {
            return <span key={i}>{val}</span>;
          }
          return (
            <a
              key={i}
              href="#"
              onClick={(e) => {
                // preact-router doesn't handle those links automatically
                e.preventDefault();
                route(`${splitPath.slice(0, i + 1).join("/")}/`)
              }}
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

export default Path;

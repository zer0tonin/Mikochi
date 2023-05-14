import { h } from "preact";
import { useState, useEffect, useContext } from "preact/hooks";
import { route } from "preact-router";

import CopyLink from "../../components/copylink";
import Download from "../../components/download";
import Icon from "../../components/icon";
// The header is directly included here to facilitate merging data from the search bar and path
import Header from "../../components/header";
import { DoubleDotPath, Path } from "../../components/path";

import { AuthContext } from "../../utils/jwt";

function formatFileSize(bytes) {
  if (bytes === 0) return "0 bytes";
  const k = 1024;
  const sizes = ["bytes", "KB", "MB", "GB", "TB", "PB", "WTF?"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  return `${size} ${sizes[i]}`;
}

const Directory = ({ dirPath = "" }) => {
  console.log(dirPath)
  const { jwt } = useContext(AuthContext);
  const [isRoot, setIsRoot] = useState(true);
  const [fileInfos, setFileInfos] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    if (dirPath != "" && !window.location.href.endsWith("/")) {
      route(`/${dirPath}/`, true);
    }
    document.title = `Mikochi ${dirPath == "" ? "" : `- /${dirPath}/`}`;
  }, [dirPath]);

  useEffect(() => {
    const fetchData = async () => {
      const params = new URLSearchParams();
      if (searchQuery != "") {
        params.append("search", searchQuery);
      }

      const response = await fetch(
        `/api/browse/${dirPath}?${params.toString()}`,
        {
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${jwt}`,
          },
        }
      );
      const json = await response.json();

      setIsRoot(json["isRoot"]);
      setFileInfos(json["fileInfos"]);
    };

    fetchData();
  }, [dirPath, searchQuery]);

  return (
    <>
      <Header searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
      <main>
        <table>
          <thead>
            <tr>
              <th />
              <th>Name</th>
              <th>Size</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {!isRoot && (
              <tr>
                <td>
                  <Icon name="folder" />
                </td>
                <td>
                  <DoubleDotPath currentDir={dirPath} />
                </td>
                <td />
                <td />
              </tr>
            )}
            {fileInfos.map((fileInfo, i) => {
              const filePath = `${dirPath == "" ? "" : `/${dirPath}`}/${
                fileInfo.path
              }`;
              if (fileInfo.isDir) {
                return (
                  <tr key={i}>
                    <td>
                      <Icon name="folder" />
                    </td>
                    <td>
                      <Path fileInfo={fileInfo} currentDir={dirPath} />
                    </td>
                    <td />
                    <td>
                      <Icon name="arrow-right-o" />
                    </td>
                  </tr>
                );
              }
              return (
                <tr key={i}>
                  <td>
                    <Icon name="file" />
                  </td>
                  <td>
                    <Path fileInfo={fileInfo} currentDir={dirPath} />
                  </td>
                  <td>{formatFileSize(fileInfo.size)}</td>
                  <td>
                    <Download filePath={filePath} />
                    <CopyLink filePath={filePath} />
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </main>
    </>
  );
};

export default Directory;

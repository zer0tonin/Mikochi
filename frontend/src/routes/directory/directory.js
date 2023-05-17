import { h } from "preact";
import { useState, useEffect, useContext } from "preact/hooks";
import { route } from "preact-router";

import CopyLink from "../../components/copylink";
import Download from "../../components/download";
import Rename from "../../components/rename";
import Delete from "../../components/delete";
import Icon from "../../components/icon";
// The header is directly included here to facilitate merging data from the search bar and path
import Header from "../../components/header";
import { DoubleDotPath, Path } from "../../components/path";
import { NameHeader, SizeHeader } from "../../components/sorting";

import { AuthContext } from "../../utils/jwt";

const formatFileSize = (bytes) => {
  if (bytes === 0) return "0 bytes";
  const k = 1024;
  const sizes = ["bytes", "KB", "MB", "GB", "TB", "PB", "WTF?"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  return `${size} ${sizes[i]}`;
}

const sorting = {
  "name_asc": (a, b) => a.path > b.path,
  "name_desc": (a, b) => a.path < b.path,
  "size_asc": (a, b) => a.size > b.size,
  "size_desc": (a, b) => a.size < b.size,
}

const Directory = ({ dirPath = "" }) => {
  const { jwt } = useContext(AuthContext);
  const [isRoot, setIsRoot] = useState(true);
  const [fileInfos, setFileInfos] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [compare, setCompare] = useState("name_asc");
  const [refresh, setRefresh] = useState(0);

  useEffect(() => {
    if (dirPath != "" && !window.location.href.endsWith("/")) {
      route(`/${dirPath}/`, true);
    }
    document.title = `Mikochi ${dirPath == "" ? "" : `- /${dirPath}/`}`;
    setSearchQuery("")

    const fetchData = async () => {
      const response = await fetch(
        `/api/browse/${dirPath}`,
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
  }, [dirPath]);

  // this two useEffect hooks look similar, but trying to combine them will get you into a race condition hell
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
  }, [searchQuery, refresh]);

  return (
    <>
      <Header searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
      <main>
        <table>
          <thead>
            <tr>
              <th />
              <NameHeader compare={compare} setCompare={setCompare} />
              <SizeHeader compare={compare} setCompare={setCompare} />
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
            {fileInfos.sort(sorting[compare]).map((fileInfo, i) => {
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
                    <td/>
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
                    <Rename
                      filePath={filePath}
                      refresh={() => setRefresh(refresh + 1)} // very hacky
                    />
                    <Delete
                      filePath={filePath}
                      refresh={() => setRefresh(refresh + 1)}
                    />
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

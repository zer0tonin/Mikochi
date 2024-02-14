import { h } from "preact";
import { useState, useEffect, useContext } from "preact/hooks";
import { useLocation } from "preact-iso";

import CopyLink from "../../components/copylink";
import Download from "../../components/download";
import Rename from "../../components/rename";
import Delete from "../../components/delete";
import Add from "../../components/add";
import Icon from "../../components/icon";

import Header from "../../components/header";
import { DoubleDotPath, Path } from "../../components/path";
import { NameHeader, SizeHeader, sorting } from "../../components/sorting";

import { AuthContext } from "../../jwt";

import "./style.css";


const formatFileSize = (bytes) => {
  if (bytes === 0) return "0 bytes";
  const k = 1024;
  const sizes = ["bytes", "KB", "MB", "GB", "TB", "PB", "WTF?"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  return `${size} ${sizes[i]}`;
};

const Directory = () => {
  const location = useLocation();
  const { jwt } = useContext(AuthContext);
  const [isRoot, setIsRoot] = useState(true);
  const [fileInfos, setFileInfos] = useState([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [compare, setCompare] = useState("name_asc");
  const [refresh, setRefresh] = useState(0); // super hacky way to trigger effects

  // we listen to the refresh event to properly handle both location changes and searches without race conditions
  useEffect(() => {
    const fetchData = async () => {
      const params = new URLSearchParams();
      if (searchQuery != "") {
        params.append("search", searchQuery);
      }

      const response = await fetch(
        `/api/browse${location.path}?${params.toString()}`,
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

    if (refresh == 0 && searchQuery == "") {
      return;
    }
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchQuery, refresh]);

  useEffect(() => {
    if (searchQuery != "") {
      setRefresh(refresh + 1);
      setCompare("none");
    }
  }, [searchQuery]);

  useEffect(() => {
    if (!location.path.endsWith("/")) {
      location.route(`${location.path}/`);
    }
    if (location.path == '/') {
      document.title = `Mikochi`;
    } else {
      document.title = `Mikochi - ${location.path}/`;
    }

    if (searchQuery != "") {
      setSearchQuery("");
      setCompare("name_asc");
    }

    setRefresh(refresh + 1);
  }, [location.path])



  return (
    <>
      <Header
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        setSearchQuery={(search) => {
          setSearchQuery(search);
          setCompare("none");
          setRefresh(refresh+1);
        }}
      />
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
                  <DoubleDotPath currentDir={location.path} />
                </td>
                <td />
                <td />
              </tr>
            )}
            {fileInfos.sort(sorting[compare]).map((fileInfo, i) => {
              let filePath;
              if (searchQuery == "" && location.path != '/') {
                filePath = `${location.path}/${fileInfo.path}`;
              } else {
                filePath = fileInfo.path;
              }
              if (fileInfo.isDir) {
                return (
                  <tr key={i}>
                    <td>
                      <Icon name="folder" />
                    </td>
                    <td>
                      <Path fileInfo={fileInfo} currentDir={location.path} />
                    </td>
                    <td />
                    <td>
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
                    <Path fileInfo={fileInfo} currentDir={location.path} />
                  </td>
                  <td>{formatFileSize(fileInfo.size)}</td>
                  <td>
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
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        <Add dirPath={location.path} refresh={refresh} setRefresh={setRefresh} />
      </main>
    </>
  );
};

export default Directory;

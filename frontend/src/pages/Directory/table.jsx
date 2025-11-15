import { h } from "preact";
import { useContext, useEffect, useState } from "preact/hooks";
import { AuthContext } from "../../jwt";
import { refresh } from ".";
import { useLocation } from "preact-iso";
import { DirActions, FileActions } from "../../components/actions";
import {
  DoubleDotFolderIcon,
  DoubleDotPath,
  FolderIcon,
  Path,
} from "../../components/path";
import Icon from "../../components/icon";
import { NameHeader, SizeHeader, sorting } from "../../components/sorting";

const formatFileSize = (bytes) => {
  if (bytes === 0) return "0 bytes";
  const k = 1024;
  const sizes = ["bytes", "KB", "MB", "GB", "TB", "PB", "WTF?"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  return `${size} ${sizes[i]}`;
};

const DirectoryTable = ({ searchQuery, compare, setCompare }) => {
  const { jwt } = useContext(AuthContext);
  const location = useLocation();
  const [isRoot, setIsRoot] = useState(true);
  const [fileInfos, setFileInfos] = useState([]);

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
        },
      );
      const json = await response.json();

      setIsRoot(json["isRoot"]);
      setFileInfos(json["fileInfos"]);
    };
    console.log("here");

    if (refresh.value == 0 && searchQuery == "") {
      return;
    }
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [refresh.value]);

  return (
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
              <DoubleDotFolderIcon currentDir={location.path} />
            </td>
            <td>
              <DoubleDotPath currentDir={location.path} />
            </td>
            <td />
            <td />
          </tr>
        )}
        {fileInfos.sort(sorting[compare]).map((fileInfo, i) => {
          let filePath = fileInfo.path;
          if (fileInfo.isDir) {
            return (
              <tr key={i}>
                <td>
                  <FolderIcon fileInfo={fileInfo} />
                </td>
                <td>
                  <Path
                    fileInfo={fileInfo}
                    currentDir={location.path}
                    isSearch={searchQuery !== ""}
                  />
                </td>
                <td />
                <td>
                  <DirActions filePath={filePath} />
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
                <Path
                  fileInfo={fileInfo}
                  currentDir={location.path}
                  isSearch={searchQuery !== ""}
                />
              </td>
              <td>{formatFileSize(fileInfo.size)}</td>
              <td>
                <FileActions filePath={filePath} />
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};

export default DirectoryTable;

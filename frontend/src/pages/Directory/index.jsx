import { h } from "preact";
import { useState, useEffect } from "preact/hooks";
import { useLocation } from "preact-iso";

import Add from "../../components/add";

import Header from "../../components/header";

import "./style.css";
import Upload from "../../components/add/upload";
import Mkdir from "../../components/add/mkdir";
import Rename from "../../components/actions/rename";
import Delete from "../../components/actions/delete";
import { signal } from "@preact/signals";
import DirectoryTable from "./table";

export const refresh = signal(0);

const Directory = () => {
  const location = useLocation();
  const [searchQuery, setSearchQuery] = useState("");
  const [compare, setCompare] = useState("name_asc");

  useEffect(() => {
    refresh.value = refresh.value + 1;
    if (searchQuery != "") {
      setCompare("none");
    } else if (compare === "none") {
      setCompare("name_asc");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchQuery]);

  useEffect(() => {
    if (!location.url.endsWith("/")) {
      location.route(`${location.path}/`, true);
    }
    if (location.path == "/") {
      document.title = `Mikochi`;
    } else {
      document.title = `Mikochi - ${decodeURI(location.path)}/`;
    }

    if (searchQuery != "") {
      setSearchQuery("");
      setCompare("name_asc");
    }

    refresh.value = refresh.value + 1;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.path]);

  return (
    <>
      <Header searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
      <main>
        <DirectoryTable
          searchQuery={searchQuery}
          compare={compare}
          setCompare={setCompare}
        />
        <Add />
        <Upload dirPath={location.path} />
        <Mkdir dirPath={location.path} />
        <Rename />
        <Delete />
      </main>
    </>
  );
};

export default Directory;

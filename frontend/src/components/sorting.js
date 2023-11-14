import { h } from "preact";

export const NameHeader = ({ compare, setCompare }) => {
  return (
    <th
      onClick={(e) => {
        e.preventDefault();
        if (compare === "name_asc") {
          setCompare("name_desc");
        } else {
          setCompare("name_asc");
        }
      }}
      style={{ cursor: "pointer" }}
    >
      Name &nbsp;
      {compare === "name_asc" && (
        <i class="gg-chevron-down" style={{ display: "inline" }} />
      )}
      {compare === "name_desc" && (
        <i class="gg-chevron-up" style={{ display: "inline" }} />
      )}
    </th>
  );
};

export const SizeHeader = ({ compare, setCompare }) => {
  return (
    <th
      onClick={(e) => {
        e.preventDefault();
        if (compare === "size_asc") {
          setCompare("size_desc");
        } else {
          setCompare("size_asc");
        }
      }}
      style={{ cursor: "pointer" }}
    >
      Size &nbsp;
      {compare === "size_asc" && (
        <i class="gg-chevron-down" style={{ display: "inline" }} />
      )}
      {compare === "size_desc" && (
        <i class="gg-chevron-up" style={{ display: "inline" }} />
      )}
    </th>
  );
};

export const sorting = {
  name_asc: (a, b) => {
    if (a.path > b.path) {
      return 1;
    }
    if (a.path == b.path) {
      return 0;
    }
    return -1;
  },
  name_desc: (a, b) => {
    if (a.path < b.path) {
      return 1;
    }
    if (a.path == b.path) {
      return 0;
    }
    return -1;
  },
  size_asc: (a, b) => {
    if (a.size > b.size) {
      return 1;
    }
    if (a.size == b.size) {
      return 0;
    }
    return -1;
  },
  size_desc: (a, b) => {
    if (a.size < b.size) {
      return 1;
    }
    if (a.size == b.size) {
      return 0;
    }
    return -1;
  },
  none: () => 0,
};

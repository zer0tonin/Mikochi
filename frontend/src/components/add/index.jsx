import { h } from "preact";
import { useEffect, useRef, useState } from "preact/hooks";

import "./style.css";
import { BigIcon } from "../icon";
import { uploadOpen } from "./upload";
import { mkdirOpen } from "./mkdir";
import { a } from "../../../dist/assets/index-BBmHDNJj";

const Mkdir = () => {
  return (
    <div class="floatingMkdir">
      <BigIcon
        name="folder-add"
        onClick={() => (mkdirOpen.value = true)}
        title="Create directory"
      />
    </div>
  );
};

const Upload = () => {
  return (
    <div class="floatingUpload">
      <BigIcon
        name="software-upload"
        onClick={() => (uploadOpen.value = true)}
        title="Upload"
      />
    </div>
  );
};

const Add = () => {
  const [extend, setExtend] = useState(false);
  const addRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (addRef.current && !addRef.current.contains(event.target)) {
        setExtend(false);
      }
    };

    if (extend) {
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [extend]);

  if (extend) {
    return (
      <div class="floating" ref={addRef}>
        <Upload />
        <Mkdir />
        <BigIcon
          name="close-o"
          onClick={() => setExtend(false)}
          title="Cancel"
        />
      </div>
    );
  }

  return (
    <div class="floating" ref={addRef}>
      <BigIcon name="add" onClick={() => setExtend(true)} title="Create" />
    </div>
  );
};

export default Add;

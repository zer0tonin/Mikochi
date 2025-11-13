import { h } from "preact";
import { useState } from "preact/hooks";

import "./style.css";
import { BigIcon } from "../icon";

const Mkdir = () => {
  return (
    <div class="floatingMkdir">
      <BigIcon
        name="folder-add"
        onClick={() => setModalOpen(true)}
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
        onClick={() => setModalOpen(true)}
        title="Upload"
      />
    </div>
  );
};

const Add = () => {
  const [extend, setExtend] = useState(false);

  if (extend) {
    return (
      <div class="floating">
        <Upload/>
        <Mkdir/>
        <BigIcon
          name="close-o"
          onClick={() => setExtend(false)}
          title="Cancel"
        />
      </div>
    );
  }

  return (
    <div class="floating">
      <BigIcon name="add" onClick={() => setExtend(true)} title="Create" />
    </div>
  );
};

export default Add;

import { h } from "preact";
import { useState } from "preact/hooks";

import "./style.css";
import Upload from "./upload";
import Mkdir from "./mkdir";
import { BigIcon } from "../icon";

const Add = ({ dirPath, refresh, setRefresh }) => {
  const [extend, setExtend] = useState(false);

  if (extend) {
    return (
      <div class="floating">
        <Upload dirPath={dirPath} refresh={refresh} setRefresh={setRefresh} />
        <Mkdir dirPath={dirPath} refresh={refresh} setRefresh={setRefresh} />
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

import { h } from "preact";
import { useState } from "preact/hooks";

import style from "./style.css";
import Upload from "./upload";
import Mkdir from "./mkdir";
import { BigIcon } from "../icon";

const Add = ({ dirPath, refresh }) => {
  const [extend, setExtend] = useState(false);

  if (extend) {
    return (
      <div class={style.floating}>
        <Upload dirPath={dirPath} refresh={refresh} />
        <Mkdir dirPath={dirPath} refresh={refresh} />
        <BigIcon
          name="close-o"
          onClick={() => setExtend(false)}
          title="Cancel"
        />
      </div>
    );
  }

  return (
    <div class={style.floating}>
      <BigIcon name="add" onClick={() => setExtend(true)} title="Create" />
    </div>
  );
};

export default Add;

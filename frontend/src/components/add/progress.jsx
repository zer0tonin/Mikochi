import { h } from "preact";

import "./style.css";

const ProgressBar = ({ done, toDo }) => {
  const progress = (done / toDo) * 100;
  return (
    <div class="progress-bar-container">
      <div class="progress-bar" style={{ width: `${progress}%` }}>
        {progress}%
      </div>
    </div>
  );
};

export default ProgressBar;

import { h } from "preact";

import "./style.css";

const ProgressBar = ({ progress }) => {
  return (
    <div class="progress-bar-container">
      <div class="progress-bar" style={{ width: `${progress}%` }}>
        {progress}%
      </div>
    </div>
  );
}

export default ProgressBar;

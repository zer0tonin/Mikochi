import { h } from "preact";

const style = {
  position: "fixed",
  bottom: "10%",
  left: "50%",
  transform: "translate(-50%, 0px)",
  "z-index": 9999,
  "background-color": "#002b36",
  "box-shadow": "0 0 5px rgba(0,0,0,.5)",
  padding: "1em",
}

const Toast = ({text, isVisible}) => {
  if (!isVisible) {
    return null;
  }

  return (
    <div style={style}>
      {text}
    </div>
  );
}

export default Toast

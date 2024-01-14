import { h } from "preact";
import { useEffect, useRef } from "preact/hooks";

import "./style.css";
import Icon from "../icon";

export const ModalHeader = ({ close, children }) => {
  return (
    <div class="modalHeader">
      {children}
      <Icon name="close" onClick={close} title="Close" />
    </div>
  );
};

export const ModalContent = ({ children }) => {
  return <div class="modalContent">{children}</div>;
};

const Modal = ({ isOpen, close, children }) => {
  const ref = useRef();
  useEffect(() => {
    const handleClickOutside = (e) => {
      if (ref.current && !ref.current.contains(e.target)) {
        close();
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [ref, close]);

  if (!isOpen) {
    return null;
  }

  return (
    <div class="modal" ref={ref}>
      {children}
    </div>
  );
};

export default Modal;

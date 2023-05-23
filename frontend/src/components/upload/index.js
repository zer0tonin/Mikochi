import { h } from "preact";
import { useContext, useEffect, useRef, useState } from "preact/hooks";

import { AuthContext } from "../../utils/jwt";
import style from "./style.css";
import Icon from "../icon";

const UploadModal = ({ isOpen, close, dirPath, refresh }) => {
  if (!isOpen) {
    return null;
  }

  const { jwt } = useContext(AuthContext);
  const [selectedFile, setSelectedFile] = useState(null);
  const [error, setError] = useState("");

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
  }, [ref]);

  const onSubmit = (e) => {
    const upload = async () => {
      const formData = new FormData();
      formData.append("file", selectedFile);

      const response = await fetch(
        dirPath != ""
          ? `/api/upload/${dirPath}/${selectedFile.name}`
          : `/api/upload/${selectedFile.name}`,
        {
          method: "PUT",
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${jwt}`,
          },
          body: formData,
        }
      );

      if (response.status !== 200) {
        setError(json["err"]);
        return;
      }
      setError("");
      close();
      refresh();
      setSelectedFile(null);
    };

    e.preventDefault();
    upload();
  };

  return (
    <div class={style.modal} ref={ref}>
      <div class={style.modalHeader}>
        File upload
        <Icon name="close" onClick={close} />
      </div>
      <div class={style.modalContent}>
        <form onSubmit={onSubmit}>
          <label class={style.fileUpload}>
            <Icon name="file-add" />
            <input
              type="file"
              class={style.input}
              onChange={(e) => setSelectedFile(e.target.files[0])}
            />
            &nbsp;
            {selectedFile != null ? (
              <span>{selectedFile.name}</span>
            ) : (
              <span>File</span>
            )}
          </label>
          <button type="submit" class={style.submit}>
            Upload
          </button>
          {error !== "" && <div class={style.error}>{error}</div>}
        </form>
      </div>
    </div>
  );
};

const Upload = ({ dirPath, refresh }) => {
  const [modalOpen, setModalOpen] = useState(false);
  return (
    <>
      <div class={style.floating}>
        <Icon name="software-upload" onClick={() => setModalOpen(true)} />
      </div>
      <UploadModal
        isOpen={modalOpen}
        close={() => setModalOpen(false)}
        dirPath={dirPath}
        refresh={refresh}
      />
    </>
  );
};

export default Upload;

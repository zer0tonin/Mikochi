import { h } from "preact";
import { useContext, useEffect, useRef, useState } from "preact/hooks";

import { AuthContext } from "../../utils/jwt";
import style from "./style.css";
import Icon, { BigIcon } from "../icon";
import Modal, { ModalContent, ModalHeader } from "../modal";

const UploadModal = ({ isOpen, close, dirPath, refresh }) => {
  if (!isOpen) {
    return null;
  }

  const { jwt } = useContext(AuthContext);
  const [selectedFile, setSelectedFile] = useState(null);
  const [error, setError] = useState("");

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
    <Modal isOpen={isOpen} close={close}>
      <ModalHeader close={close}>File upload</ModalHeader>
      <ModalContent>
        <form onSubmit={onSubmit}>
          <label class={style.fileUpload}>
            <Icon name="file-add" />
            <input
              type="file"
              class={style.input}
              onChange={(e) => setSelectedFile(e.target.files[0])}
              aria-label="Select a file"
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
      </ModalContent>
    </Modal>
  );
};

const Upload = ({ dirPath, refresh }) => {
  const [modalOpen, setModalOpen] = useState(false);
  return (
    <>
      <div class={style.floatingUpload}>
        <BigIcon
          name="software-upload"
          onClick={() => setModalOpen(true)}
          title="Upload"
        />
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

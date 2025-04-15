import { h } from "preact";
import { useContext, useState } from "preact/hooks";

import { AuthContext } from "../../jwt";
import Icon, { BigIcon } from "../icon";
import Modal, { ModalContent, ModalHeader } from "../modal";
import "./style.css";
import handleError, { getHttpErrorDescription } from "../../error";
import Toast from "../toast";

const UploadModal = ({
  isOpen,
  close,
  dirPath,
  refresh,
  setRefresh,
  setSuccess,
}) => {
  const { jwt } = useContext(AuthContext);
  const [selectedFile, setSelectedFile] = useState(null);
  const [error, setError] = useState("");

  if (!isOpen) {
    return null;
  }

  const onSubmit = (e) => {
    e.preventDefault();

    // We use XHR instead of fetch because it provides a simpler way to check upload progress
    var xhr = new XMLHttpRequest();
    xhr.open(
      'PUT',
      dirPath != ""
        ? `/api/upload/${dirPath}/${selectedFile.name}`
        : `/api/upload/${selectedFile.name}`,
    )
    xhr.setRequestHeader("Accept", "application/json")
    xhr.setRequestHeader("Authorization", `Bearer ${jwt}`)

    xhr.onload = () => {
      if (xhr.status !== 200) {
        setError(getHttpErrorDescription(xhr.status));
      } else {
        setError("");
        setSelectedFile(null);
        setRefresh(refresh + 1);
        close();
        setSuccess(true);
        setTimeout(() => {
          setSuccess(false);
        }, 2000);
      }
    }
    xhr.onerror = () => {
      setError(getHttpErrorDescription(xhr.status));
    }

    const formData = new FormData();
    formData.append("file", selectedFile);

    xhr.send(formData);
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <ModalHeader close={close}>File upload</ModalHeader>
      <ModalContent>
        <form onSubmit={onSubmit}>
          <label class="fileUpload">
            <Icon name="file-add" />
            <input
              type="file"
              class="hiddenInput"
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
          <button type="submit" class="submit">
            Upload
          </button>
          {error !== "" && <div class="error">{error}</div>}
        </form>
      </ModalContent>
    </Modal>
  );
};

const Upload = ({ dirPath, refresh, setRefresh }) => {
  const [modalOpen, setModalOpen] = useState(false);
  const [success, setSuccess] = useState(false);
  return (
    <>
      <div class="floatingUpload">
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
        setRefresh={setRefresh}
        setSuccess={setSuccess}
      />
      {success && (
        <Toast text="File Uploaded Successfully" isVisible={success} />
      )}
    </>
  );
};

export default Upload;

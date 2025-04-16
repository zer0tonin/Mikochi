import { h } from "preact";
import { useContext, useState } from "preact/hooks";

import { AuthContext } from "../../jwt";
import Icon, { BigIcon } from "../icon";
import Modal, { ModalContent, ModalHeader } from "../modal";
import "./style.css";
import { getHttpErrorDescription } from "../../error";
import Toast from "../toast";
import ProgressBar from "./progress";

const UploadModal = ({
  isOpen,
  close,
  dirPath,
  refresh,
  setRefresh,
  setSuccess,
}) => {
  const { jwt } = useContext(AuthContext);
  const [selectedFiles, setSelectedFiles] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState("");

  if (!isOpen) {
    return null;
  }

  const onSubmit = (e) => {
    e.preventDefault();

    const upload = (file) => {
      // We use XHR instead of fetch because it provides a simpler way to check upload progress
      var xhr = new XMLHttpRequest();
      xhr.open(
        'PUT',
        dirPath != ""
          ? `/api/upload/${dirPath}/${file.name}`
          : `/api/upload/${file.name}`,
      )
      xhr.setRequestHeader("Accept", "application/json")
      xhr.setRequestHeader("Authorization", `Bearer ${jwt}`)

      xhr.onload = () => {
        setUploading(false);
        if (xhr.status !== 200) {
          setError(getHttpErrorDescription(xhr.status));
        } else {
          setError("");
          setSelectedFiles(null);
          setRefresh(refresh + 1);
          close();
          setSuccess(true);
          setTimeout(() => {
            setSuccess(false);
          }, 2000);
        }
      }

      xhr.onerror = () => {
        setUploading(false);
        setError(getHttpErrorDescription(xhr.status));
      }

      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          setProgress((event.loaded / event.total) * 100)
        }
      }

      const formData = new FormData();
      formData.append("file", file);
      xhr.send(formData);

      setUploading(true);
    }

    selectedFiles.map(upload)
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <ModalHeader close={close}>File upload</ModalHeader>
      <ModalContent>
        { uploading ?
          <ProgressBar progress={progress} />
          : (
          <form onSubmit={onSubmit}>
            <label class="fileUpload">
              <input
                type="file"
                class="hiddenInput"
                onChange={(e) => setSelectedFiles(Array.from(e.target.files))}
                aria-label="Select a file"
                multiple
              />
              {selectedFiles != null ? (
                selectedFiles.map((f) => (
                  <div>
                    <Icon name="file-add" />
                    &nbsp;
                    <span>{f.name}</span>
                  </div>
                ))
              ) : (
                <>
                  <Icon name="file-add" />
                  &nbsp;
                  <span>File</span>
                </>
              )}
            </label>
            <button type="submit" class="submit" disabled={selectedFiles == null}>
              Upload
            </button>
            {error !== "" && <div class="error">{error}</div>}
          </form>
        )}
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

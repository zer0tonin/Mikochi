import { h } from "preact";
import { useContext, useState } from "preact/hooks";
import { signal } from "@preact/signals";

import { AuthContext } from "../../jwt";
import Icon from "../icon";
import Modal, { ModalContent, ModalHeader } from "../modal";
import "./style.css";
import { getHttpErrorDescription } from "../../error";
import Toast from "../toast";
import ProgressBar from "./progress";

export const uploadOpen = signal(false);

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
  const [uploaded, setUploaded] = useState(0);
  const [uploadSize, setUploadSize] = useState(0);
  const [error, setError] = useState("");

  if (!isOpen) {
    return null;
  }

  const onSubmit = (e) => {
    e.preventDefault();

    const upload = (file) => {
      return new Promise((resolve, reject) => {
        // We use XHR instead of fetch because it provides a simpler way to check upload progress
        var xhr = new XMLHttpRequest();
        xhr.open(
          "PUT",
          dirPath != ""
            ? `/api/upload/${dirPath}/${file.name}`
            : `/api/upload/${file.name}`,
        );
        xhr.setRequestHeader("Accept", "application/json");
        xhr.setRequestHeader("Authorization", `Bearer ${jwt}`);

        xhr.onload = () => {
          if (xhr.status !== 200) {
            reject(getHttpErrorDescription(xhr.status));
          } else {
            resolve();
          }
        };

        xhr.onerror = () => {
          reject(getHttpErrorDescription(xhr.status));
        };

        xhr.upload.onloadstart = (event) => {
          if (event.lengthComputable) {
            setUploadSize(uploadSize + event.total);
          }
        };

        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable) {
            setUploaded(uploaded + event.loaded);
          }
        };

        const formData = new FormData();
        formData.append("file", file);
        xhr.send(formData);

        setUploading(true);
      });
    };

    Promise.all(selectedFiles.map(upload))
      .then(() => {
        setError("");
        setSelectedFiles(null);

        setRefresh(refresh + 1);
        setSuccess(true);
        setTimeout(() => {
          setSuccess(false);
        }, 2000);
        close();
      })
      .catch((err) => {
        setError(err);
      })
      .finally(() => {
        setUploading(false);
        setUploaded(0);
        setUploadSize(0);
      });
  };

  return (
    <Modal isOpen={isOpen} close={close}>
      <ModalHeader close={close}>File upload</ModalHeader>
      <ModalContent>
        {uploading ? (
          <ProgressBar done={uploaded} toDo={uploadSize} />
        ) : (
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
            <button
              type="submit"
              class="submit"
              disabled={selectedFiles == null}
            >
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
  const [success, setSuccess] = useState(false);
  return (
    <>
      <UploadModal
        isOpen={uploadOpen.value}
        close={() => (uploadOpen.value = false)}
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

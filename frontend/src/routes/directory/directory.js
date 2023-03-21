import { h } from 'preact';
import { useState, useEffect } from 'preact/hooks';

import Icon from '../../components/icon';


function formatFileSize(bytes) {
  if (bytes === 0) return '0 bytes';
  const k = 1024;
  const sizes = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'WTF?'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  return `${size} ${sizes[i]}`;
}


const Directory = ({ dirPath = '', search = '' }) => {
    const [isRoot, setIsRoot] = useState(true)
    const [fileInfos, setFileInfos] = useState([])

    useEffect(() => {
        const fetchData = async () => {
            let json = null;
            if (search == '') {
                const response = await fetch(`/api/browse/${dirPath}`);
                json = await response.json();
            } else {
                const response = await fetch(`/api/search/${search}`);
                json = await response.json();
            }
            
            setIsRoot(json['isRoot'])
            setFileInfos(json['fileInfos'])
        };

        fetchData();
    }, [dirPath, search]);

    return (
        <table>
            <thead>
                <tr>
                    <th></th>
                    <th>Name</th>
                    <th>Size</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                { !isRoot &&
                    <tr>
                        <td><Icon name="folder" /></td>
                        <td><a href="..">..</a></td>
                        <td></td>
                        <td></td>
                    </tr>
                }
                { fileInfos.map((fileInfo) => {
                    if (fileInfo.isDir) {
                        return (
                            <tr>
                                <td><Icon name="folder" /></td>
                                <td><a href={`${fileInfo.name}/`}>{fileInfo.name}</a></td>
                                <td></td>
                                <td><Icon name="arrow-right-o" /></td>
                            </tr>
                        )
                    }
                    console.log(fileInfo)
                    return (
                        <tr>
                            <td><Icon name="file" /></td>
                            <td>{fileInfo.name}</td>
                            <td>{formatFileSize(fileInfo.size)}</td>
                            <td><Icon name="arrow-down-o" /><Icon name="copy" /></td>
                        </tr>
                    );
                })}
            </tbody>
        </table>
    )
}

export default Directory

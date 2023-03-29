import { h } from 'preact';
import { useState, useEffect } from 'preact/hooks';
import { route } from 'preact-router';

import CopyLink from '../../components/copylink'
import Icon from '../../components/icon';
// The header is directly included here to facilitate merging data from the search bar and path
import Header from '../../components/header';
import Path from '../../components/path';


function formatFileSize(bytes) {
  if (bytes === 0) return '0 bytes';
  const k = 1024;
  const sizes = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'WTF?'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  return `${size} ${sizes[i]}`;
}


const Directory = ({ dirPath = '' }) => {
    if (dirPath != '' && !window.location.href.endsWith('/')) {
        route(`/${dirPath}/`, true)
    }

    const [isRoot, setIsRoot] = useState(true)
    const [fileInfos, setFileInfos] = useState([])
    const [searchQuery, setSearchQuery] = useState('')

    useEffect(() => {
        const fetchData = async () => {
            const params = new URLSearchParams()
            if (searchQuery != '') {
                params.append('search', searchQuery)
            }

            
            const response = await fetch(`/api/browse/${dirPath}?${params.toString()}`);
            const json = await response.json();
            
            setIsRoot(json['isRoot'])
            setFileInfos(json['fileInfos'])
        };

        fetchData();
    }, [dirPath, searchQuery]);

    return (
        <>
            <Header searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
            <main>
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
                                        <td><Path fileInfo={fileInfo} /></td>
                                        <td></td>
                                        <td><Icon name="arrow-right-o" /></td>
                                    </tr>
                                )
                            }
                            return (
                                <tr>
                                    <td><Icon name="file" /></td>
                                    <td><Path fileInfo={fileInfo} /></td>
                                    <td>{formatFileSize(fileInfo.size)}</td>
                                    <td><Icon name="arrow-down-o" /><CopyLink filePath={`${dirPath == "" ? "" : "/" + dirPath}/${fileInfo.path}`} /></td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            </main>
        </>
    )
}

export default Directory

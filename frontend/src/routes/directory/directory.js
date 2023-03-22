import { h } from 'preact';
import { useState, useEffect } from 'preact/hooks';

import Icon from '../../components/icon';
// The header is directly included here to facilitate merging data from the search bar and path
import Header from '../../components/header';


function formatFileSize(bytes) {
  if (bytes === 0) return '0 bytes';
  const k = 1024;
  const sizes = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'WTF?'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const size = parseFloat((bytes / Math.pow(k, i)).toFixed(2));
  return `${size} ${sizes[i]}`;
}


const Directory = ({ dirPath = '' }) => {
    const [isRoot, setIsRoot] = useState(true)
    const [fileInfos, setFileInfos] = useState([])
    const [searchQuery, setSearchQuery] = useState('')

    useEffect(() => {
        const fetchData = async () => {
            const response = await fetch(`/api/browse/${dirPath}`);
            const json = await response.json();
            
            setIsRoot(json['isRoot'])
            setFileInfos(json['fileInfos'])
        };

        fetchData();
    }, [dirPath]);

    useEffect(() => {
        const fetchData = async () => {
            if (searchQuery == '') {
                return
            }

            const params = new URLSearchParams({ search: searchQuery })
            const response = await fetch(`/api/browse/${dirPath}?${params.toString()}`);
            const json = await response.json();
            
            setFileInfos(json['fileInfos'])
        };

        // we use setTimeout to wait until the user stops typing before making a query
        const timer = setTimeout(fetchData, 500)
        return () => clearTimeout(timer)
    }, [searchQuery]);

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
                                        <td><a href={`${fileInfo.path}/`}>{fileInfo.path}</a></td>
                                        <td></td>
                                        <td><Icon name="arrow-right-o" /></td>
                                    </tr>
                                )
                            }
                            return (
                                <tr>
                                    <td><Icon name="file" /></td>
                                    <td>{fileInfo.path}</td>
                                    <td>{formatFileSize(fileInfo.size)}</td>
                                    <td><Icon name="arrow-down-o" /><Icon name="copy" /></td>
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

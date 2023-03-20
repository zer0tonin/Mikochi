import { h } from 'preact';
import { useState, useEffect } from 'preact/hooks';
import { Link } from 'preact-router/match';

import Icon from '../../components/icon';

const Directory = (dirPath) => {
    const [isRoot, setIsRoot] = useState(true)
    const [dirEntries, setDirEntries] = useState([])

    useEffect(() => {
        setIsRoot(true)
        setDirEntries([{name: "file1.mp4", isDir: false}, { name: "file2.mp4", isDir: false}])
    }, [dirPath]);

    return (
        <table class="table">
            <thead>
                <tr>
                    <th></th>
                    <th>Name</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                { !isRoot &&
                    <tr>
                        <td><Icon name="folder" /></td>
                        <td><Link href="..">..</Link></td>
                    </tr>
                }
                { dirEntries.map((dirEntry) => 
                    <tr>
                        { dirEntry.isDir
                            ? <td><Icon name="folder" /></td>
                            : <td><Icon name="file" /></td>
                        }
                        { dirEntry.isDir
                            ? <td><Link href="{{ .Name }}/">{dirEntry.name}</Link></td>
                            : <td>{dirEntry.name}</td>
                        }
                        { dirEntry.isDir
                            ? <td><Icon name="arrow-right-o" /></td>
                            : <td><Icon name="arrow-down-o" /><Icon name="copy" /></td>
                        }
                    </tr>
                )}
            </tbody>
        </table>
    )
}

export default Directory

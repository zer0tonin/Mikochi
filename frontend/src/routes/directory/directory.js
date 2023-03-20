import { h } from 'preact';
import { useState, useEffect } from 'preact/hooks';

const Directory = () => {
    const [path, setPath] = useState("/")
    const [isRoot, setIsRoot] = useState(true)
    const [dirEntries, setDirEntries] = useState([]))

    useEffect(() => {
    }

    return (
        <div>
            <table class="table">
                <thead>
                    <tr>
                        <th></th>
                        <th>Name</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    { isRoot &&
                        <tr>
                            <td><i class="gg-folder"></i></td>
                            <td><a href="..">..</a></td>
                        </tr>
                    }
                    { dirEntries.map((dirEntry) => 
                        <tr>
                            { dirEntry.isDir
                                ? <td><i class="icon gg-folder"></i></td>
                                : <td><i class="icon gg-file"></i></td>
                            }
                            { dirEntry.isDir
                                ? <td><a href="{{ .Name }}/">{{ .Name }}</a></td>
                                : <td>{{ .Name }}</td>
                            }
                            { dirEntry.isDir
                                ? <td><i class="icon gg-arrow-right-o"></i></td>
                                : <td><i class="icon gg-arrow-down-o"></i><i class="icon gg-copy"></i></td>
                            }
                        </tr>
                    )}
                    {{ end }}
                </tbody>
            </table>
        </div>
    )
}

export default Directory

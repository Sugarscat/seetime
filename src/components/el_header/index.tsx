import './index.css'
import IconGitHub from "../icon/IconGitHub";
import IconDocs from "../icon/IconDocs";

function ElHeader() {
    return(
        <div className="el-header">
            <div className={"logo"}>

            </div>
            <ul>
                <li>
                    <a href="#" target="_blank">
                        <IconDocs/>
                         文档
                    </a>
                </li>
                <li>
                    <a href="#" target="_blank"><IconGitHub/></a>
                </li>
            </ul>
        </div>
    )
}

export default ElHeader

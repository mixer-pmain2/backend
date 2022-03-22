import React, {useState} from "react";

const Card = ({title, img, onClick}) => {
    const [hover, setHover] = useState(false)
    let style = {}
    const hoverStyle = {borderTopWidth: 5, borderRightWidth: 5, marginTop: -5, marginRight: -5, transitionDuration: "100ms"}

    const toggleHover = () => {
        setHover(!hover)
    }

    style = hover ? hoverStyle : {}

    return <div className="card p-2" onClick={onClick} style={{...style}} onMouseEnter={toggleHover} onMouseLeave={toggleHover}>
        <img src={img.default} height={"80"} className="card-img-top" alt="..."/>
        <div className="card-body">
            <h5 className="card-title">{title}</h5>
            {/*<a href="#" className="btn btn-primary">Go somewhere</a>*/}
        </div>
    </div>
}

export default Card
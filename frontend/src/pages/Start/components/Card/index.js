import React, {useState} from "react";

const Card = ({title, img, onClick}) => {
    const [hover, setHover] = useState(false)
    let style = {}
    const cWidth = 250;

    const hoverStyle = {
        borderTopWidth: 5,
        borderRightWidth: 5,
        paddingBottom: 0,
        paddingLeft: 0,
        // transitionDuration: "100ms"
    }

    const toggleHover = () => {
        setHover(!hover)
    }

    style = hover ? hoverStyle : {paddingBottom: 5}
    const styleWrapper = hover ? {paddingBottom: 0} : {paddingBottom: 5}

    return <div className="col-sm-auto mb-4" style={{cursor: "pointer", ...styleWrapper}}>
        <div className="card p-2" onClick={onClick} style={{height: "100%", width: cWidth, ...style}}
             onMouseEnter={toggleHover} onMouseLeave={toggleHover}>
            <img src={img.default} height={"80"} className="card-img-top mt-3" alt="..."
                 // style={hover ? {marginTop: 0, marginRight: 0} : {marginTop: 5, marginRight: 5}}
            />
            <div className="card-body"
                 // style={hover ? {paddingBottom: 0} : {height: 70, paddingBottom: 5}}
            >
                <h5 className="card-title">{title}</h5>
            </div>
        </div>
    </div>
}

export default Card
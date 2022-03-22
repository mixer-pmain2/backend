import React from "react";

const Pagination = ({current, total, nextPage, prevPage, onPage}) => {

    const pagesArr = Array(total).fill().map((_, i) => i)
    const isFirst = current === 0
    const isLast = current === total - 1

    let beforeCurr = false;
    let afterCurr = false;

    const checkDotBtn = (v) => {
        if (v <= current && !beforeCurr) {
            beforeCurr = true
            return false
        } else if (v >= current && !afterCurr) {
            afterCurr = true
            return false
        }
        return true
    }

    return <nav>
        <ul className="pagination">
            <li className={`page-item ${isFirst && "disabled"}`} disabled={isFirst}>
                <a className="page-link" href="#" aria-label="Previous" onClick={prevPage}>
                    <span aria-hidden="true">&laquo;</span>
                </a>
            </li>
            {
                pagesArr.map(v => {
                        return (v === 0 || v === total - 1 || (v >= current - 2 && v <= current + 2)) ?
                            <li className={`page-item ${current === v ? "active" : ""}`} key={v}>
                                <a className="page-link" href="#" onClick={_ => onPage(v)}>
                                    {v + 1}
                                </a>
                            </li> :
                            (!checkDotBtn(v) && <li className={`page-item`} key={v}>
                                <a className="page-link" href="#">
                                    ...
                                </a>
                            </li>)
                    }
                )
            }
            <li className={`page-item ${isLast && "disabled"}`} disabled={isLast}>
                <a className="page-link" href="#" aria-label="Next" onClick={nextPage}>
                    <span aria-hidden="true">&raquo;</span>
                </a>
            </li>
        </ul>
    </nav>
}

export default Pagination
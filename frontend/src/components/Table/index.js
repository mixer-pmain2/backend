import React, {useEffect, useState} from "react"
import Pagination from "../Pagination";
import Loading from "../Loading";

const PAGE_SIZE = 10

const Table = ({columns, data, mapper, onDoubleClick, selecting, onClick, loading}) => {
    const [currentPage, setCurrentPage] = useState(0)
    const [selectedRow, setSelectedRow] = useState(false)
    const selectedStyle = selecting ? {backgroundColor: "gray", color: "white"} : {}

    let total = Math.floor(data.length / PAGE_SIZE)
    if (data.length % PAGE_SIZE > 0) total += 1

    const handleNextPage = () => {
        if (currentPage >= total) return
        setCurrentPage(currentPage + 1)
    }

    const handlePrevPage = () => {
        if (currentPage <= 0) return
        setCurrentPage(currentPage - 1)
    }

    const handleDoubleClick = (e) => {
        onDoubleClick && onDoubleClick(e)
    }

    const onPage = (num) =>
        setCurrentPage(num)

    const showDataOnlyForPage = (i) => currentPage * PAGE_SIZE <= i && (currentPage * PAGE_SIZE + PAGE_SIZE) > i

    const handleClick = (i) => {
        setSelectedRow(i)
        onClick && onClick(data[i])
    }

    useEffect(() => {
        setCurrentPage(0)
    }, [data])

    return <div>
        <table className="table table-striped table-hover">
            {
                columns && <thead>
                <tr>
                    {columns.map((v, i) => <th scope="col" key={i}>{v}</th>)}
                </tr>
                </thead>
            }
            <tbody>
            {
                loading ? <tr>
                    <td colSpan={columns.length}>
                        <div className="d-flex justify-content-center">
                            <Loading isLoading={true}/>
                        </div>
                    </td>
                </tr> : <>
                    {data?.map((row, i) => showDataOnlyForPage(i) &&
                        <tr key={i}
                            onDoubleClick={_ => handleDoubleClick(row)}
                            onClick={_ => handleClick(i)}
                            style={selectedRow === i ? selectedStyle : {}}
                        >
                            {mapper(row)}
                        </tr>)}
                    {data.length === 0 && <tr>
                        <td colSpan={columns.length}>
                            <div className="d-flex justify-content-center">Нет данных</div>
                        </td>
                    </tr>}</>
            }
            </tbody>

        </table>
        <div className="d-flex justify-content-end">
            {total > 1 && <Pagination
                total={total}
                current={currentPage}
                nextPage={handleNextPage}
                prevPage={handlePrevPage}
                onPage={onPage}
            />}
        </div>
    </div>
}

export default Table
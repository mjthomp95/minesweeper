import React from 'react';
import Mark from './Mark';
import Show from './Show';
import Blank from './Blank';

const Cell = (props) => {

    const clickHandler = (e) => {
        click(e.button, props.row, props.col, props.changes);
    };
    const doubleClickHandler = (e) => {
        doubleClick(props.row, props.col, props.changes);
    };

    let char = String.fromCharCode(props.value)

    return (<div className='cell' onClick={clickHandler} onDoubleClick={doubleClickHandler}
            onContextMenu={(e) => {e.preventDefault(); click(e.button, props.row, props.col, props.changes)}}>
            {   (show(props.value)) ?
                <Show value={char} /> :
                (mark(char)) ?
                <Mark /> :
                <Blank />
            }</div>);
}

const fetchChanges = (row, col, method, callBack) => {
    // TODO: stop doubleClick triggering click event.
    fetch(`http://localhost:8080/${method}`, {
        method: 'POST',
        headers: {'Content-Type':'application/x-www-form-urlencoded'},
        body: `row=${row}&col=${col}`
    })
    .then(resp => {return resp.json()}) // IDEA: Could I change so it only sends if there are multiple changes and not just this one.
    .then(data => {callBack(data)});
}

const click = (button, row, col, callBack) => {
    if (button === 0){
        fetchChanges(row, col, 'choose', callBack);
    } else if (button === 2){
        fetchChanges(row, col, 'mark', callBack);
    }

}

const doubleClick = (row, col, callBack) => {
    fetchChanges(row, col, 'expand', callBack);
}

const show = (val) => {
    if (val === 0) {
        return false;
    }
    return true;
}

const mark = (val) => {
    if (val === 'x') {
        return true;
    }
    return false;
}

export default Cell;

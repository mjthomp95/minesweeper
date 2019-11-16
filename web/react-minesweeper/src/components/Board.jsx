import React, {useState} from 'react';
import Cell from './Cell';

// TODO: evaluate for performance; tried 50x50 and froze
// IDEA: use a reducer for board
const Board = React.memo((props) => {
    const [board, setBoard] = useState([]);

    const newBoard = () => {
        getNewBoard(props.height, props.width).then(data => {
            if (data) {
                setBoard(data.cells);
                props.setMines(data.mines);
            } else {
                // TODO: deal with error
                props.setGame({gameStart: false, gameOver: false});
            }});
    };

    if (!props.start && !props.stop && board.length > 0) {
        setBoard([]);
    } else if (props.start && props.stop) {
        props.setGame({gameStart: true, gameOver: false})
        newBoard();
    }

    if (props.start && board.length === 0) {
        newBoard();
    }

    const getChanges = (changes) => {
        if (changes.error) {
            if (changes.error === 'Lose' || changes.error === 'Win') {
                //game over indicate win or lose
                //Should still have changes to board
                props.setGame({gameStart: false, gameOver: true, win: changes.error});
            } else {
                // TODO: Display error
                console.log(changes.error);
                return
            }
        }

        if (changes.mines || changes.mines === 0) {
            props.setMines(changes.mines);
        }

        setBoard(prevBoard => {
            let nextBoard = [...prevBoard];
            if (changes.cells) {
                let cell;
                for (cell of changes.cells) {
                    nextBoard[cell.Row][cell.Col] = cell.Value;
                }
            }
            return nextBoard;
        });


    };

// IDEA: timer might come from server in later version.
// <!-- transparent image with 1:1 intrinsic aspect ratio -->
    const boardStyle = {
        gridTemplateColumns: " 1fr".repeat(props.width),
        gridTemplateRows: " 1fr".repeat(props.height),
    }

    return (<div className='boardContainer'>
                <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
                alt="transparent with 1:1 intrinsic aspect ratio"/>
                <div className='board' style={boardStyle}>
                    {board && [...board].map((row, i) =>
                        row.map((cell, j) =>
                            <Cell key={i*10+j} row={i} col={j} value={cell} changes={getChanges}/>
                    ))}
                </div>
            </div>);
});
// TODO: stop the board moving when "Win" or "Lose" and when scoreboard goes away.
// IDEA: add overlay for board for win/lose
// TODO: add height and width input for other sized board

const getNewBoard = (height = 10, width = 10) => {
    const json = fetch("http://localhost:8080/new", {
        method: 'POST',
        headers: {'Content-Type':'application/x-www-form-urlencoded'},
        body: `height=${height}&width=${width}`
    })
    .then(resp => {return resp.json()})
    .then(data => {console.log(data); return data})
    .catch(err => console.log(err))
    return json
}

export default Board;

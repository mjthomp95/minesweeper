import React, {useState} from 'react';
import Cell from './Cell';
import ScoreBoard from './ScoreBoard';

const Board = () => {
    const [board, setBoard] = useState({cells: [], mines: 0, numCells: 0});
    const [gameState, setGameState] = useState({gameStart: false, gameOver: false})
    const newBoard = () => {
        if (board.cells) {
            endGame();
        }
        getNewBoard().then(data => {
            if (data) {
                setBoard(data)
                setGameState({gameStart: true, gameOver: false})
            }});
    };
    const endGame = () => {
        setGameState({gameStart: false, gameOver: false})
        fetch("http://localhost:8080/end")
    };

    const getChanges = (changes) => {
        if (changes.error) {
            if (changes.error === 'Lose' || changes.error === 'Win') {
                //game over indicate win or lose
                //Should still have changes to board
                setGameState({gameStart: false, gameOver: true, win: changes.error})
            } else {
                //display error
                console.log(changes.error)
                return
            }
        }

        let tmpBoard = {...board}
        if (changes.cells) {
            let cell;
            for (cell of changes.cells) {
                tmpBoard.cells[cell.Row][cell.Col] = cell.Value
            }
        }

        if (changes.mines || changes.mines === 0) {
            tmpBoard.mines = changes.mines
        }
        if (changes.numCells || changes.numCells === 0) {
            tmpBoard.numCells = changes.numcells
        }
        setBoard(tmpBoard)
    };

// IDEA: timer might come from server in later version.
// <!-- transparent image with 1:1 intrinsic aspect ratio -->

    return (<div className='canvas'>
            {gameState.gameOver && <h2>{gameState.win}</h2>}
            <ScoreBoard numMines={board.mines} start={gameState.gameStart} stop={gameState.gameOver}
            reset={!gameState.gameStart && !gameState.gameOver}/>
            <div className='boardContainer'>
                <img src="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
                alt="transparent with 1:1 intrinsic aspect ratio"/>
                <div className='board'>
                    {board.cells ? [...board.cells].map((row, i) =>
                        row.map((cell, j) =>
                            <Cell key={i*10+j} row={i} col={j} value={cell} changes={getChanges}/>
                    )) : <div></div>}
                </div>
            </div>
            <br/>
            <div className='panel'>
                <button onClick={newBoard}>
                    New Game
                </button>
                <button onClick={endGame}>
                    End Game
                </button>
            </div>
            </div>);
}
// TODO: stop the board moving when "Win" or "Lose" and when scoreboard goes away.
// IDEA: add overlay for board for win/lose
// TODO: add height and width input for other sized board

const getNewBoard = () => {
    const json = fetch("http://localhost:8080/new", {
        method: 'POST',
        headers: {'Content-Type':'application/x-www-form-urlencoded'},
        body: 'height=10&width=10'
    })
    .then(resp => {return resp.json()})
    .then(data => {return data})
    .catch(err => console.log(err))
    return json
}

export default Board;

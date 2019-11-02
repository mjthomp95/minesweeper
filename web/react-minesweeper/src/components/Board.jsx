import React, {useState} from 'react';
import Cell from './Cell';

const Board = () => {
    const [board, setBoard] = useState({cells: [], mines: 0, numCells: 0});
    const [gameState, setGameState] = useState({gameOver: false, win: 'Lose'})
    const newBoard = () => {
        getNewBoard().then(data => {
            if (data) {
                setBoard(data)
            }});
    };
    const endGame = () => {
        setGameState({gameOver: false})
        fetch("http://localhost:8080/end")
    };

    const getChanges = (changes) => {
        if (changes.error) {
            if (changes.error === 'Lose' || changes.error === 'Win') {
                //game over indicate win or lose
                //Should still have changes to board
                setGameState({gameOver: true, win: changes.error})
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

        if (changes.mines) {
            tmpBoard.mines = changes.mines
        }
        if (changes.numCells) {
            tmpBoard.numCells = changes.numcells
        }
        setBoard(tmpBoard)
    };

    return (<div className='canvas'>
            {gameState.gameOver && <h2>{gameState.win}</h2>}
            <div className='board'>
                {[...board.cells].map((row, i) =>
                    row.map((cell, j) =>
                        <Cell key={i*10+j} row={i} col={j} value={cell} changes={getChanges}/>
                ))}
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

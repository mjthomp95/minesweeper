import React, {useState, useRef} from 'react';
import Board from './Board';
import ScoreBoard from './ScoreBoard';

// IDEA: custom hook for gameState
const Canvas = React.memo(() => {
    const [gameState, setGameState] = useState({gameStart: false, gameOver: false});
    const [dimensions, setDimensions] = useState({height: 0, width: 0});
    const [mines, setMines] = useState(0);
    const height = useRef(null);
    const width = useRef(null);

    const newGame = () => {
        // FIXME: only works with square boards right now. can change by calculating aspect ratio then changing img tag for aspect ratio.
        if (height.current.value <= 0 || width.current.value <= 0) {
            return;
        }
        restartGame();
        setDimensions({height: height.current.value, width: width.current.value});
        setGameState(prev => {return {...prev, gameStart: true}});
    };

    const endGameClosure = (stateChange) => {return (() => {
        if (stateChange) {
            setGameState({gameStart: false, gameOver: false});
            setMines(0);
        }
        fetch("http://localhost:8080/end"); // TODO: catch error?
    })};

    const endGame = endGameClosure(true);
    const restartGame = endGameClosure(false);

    return (<div className='canvas'>
            {gameState.gameOver && <h2>{gameState.win}</h2>}
            <ScoreBoard numMines={mines} start={gameState.gameStart} stop={gameState.gameOver}
            reset={!gameState.gameStart && !gameState.gameOver}/>
            <Board setGame={setGameState} height={dimensions.height} width={dimensions.width}
            start={gameState.gameStart} stop={gameState.gameOver}
            mines={mines} setMines={setMines}/>
            <br/>
            <div className='panel'>
                {!gameState.gameStart && <form>
                <label>
                Height: <input name='height' type='number' defaultValue={10} ref={height}/><br/>
                </label>
                <label>
                Width: <input name='width' type='number' defaultValue={10} ref={width}/><br/>
                </label>
                </form>}
                <button onClick={newGame}>
                    New Game
                </button>
                <button onClick={endGame}>
                    End Game
                </button>
            </div>
            </div>)
});

export default Canvas;

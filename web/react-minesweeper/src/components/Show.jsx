import React from 'react';

const Show = React.memo((props) => {
    // TODO: make a mine graphic
    return (<div><div  className='show'>{props.value}</div></div>);
});

export default Show;

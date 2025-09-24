import {useState, useEffect} from 'react';
import {Navigate} from 'react-router-dom';

const ProtectedRoute = ({children}) => {
    const [isAuthenticated, setIsAuthenticated] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(()=>{
        checkAuth();
    },[]);

    const checkAuth = async () => {
        try {
            const response = await fetch('/api/user-status', {
            credentials: 'include'
            });
            if (response.ok) {
                setIsAuthenticated(true);
            }else{
                setIsAuthenticated(false);
            }
        }catch (error){
            console.error('Auth check failed:',error);
            setIsAuthenticated(false);
        }finally{
            setIsLoading(false);
        }
    };

    if (isLoading){
            return (
        <div style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            background: 'linear-gradient(135deg, #0f0f23, #1a1a2e)',
            color: '#00fff7',
            fontFamily: 'JetBrains Mono, monospace'
        }}>
            <div>Loading...</div>
        </div>
        ); 
    }

    if (!isAuthenticated){
         return <Navigate to="/auth" replace />;
    }

    return children;
};

export default ProtectedRoute;
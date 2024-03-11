import PropTypes from 'prop-types';

// material-ui
import { useTheme } from '@mui/material/styles';
import { Avatar, Box, Button, ButtonBase } from '@mui/material';

// project imports
import LogoSection from '../LogoSection';
import ProfileSection from './ProfileSection';

// assets
import { IconMenu2 } from '@tabler/icons-react';
import { Link } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { API } from '../../../utils/api';
import { showError } from '../../../utils/common';
import { useSelector } from 'react-redux';

// ==============================|| MAIN NAVBAR / HEADER ||============================== //

const Header = ({ handleLeftDrawerToggle }) => {
  const theme = useTheme();
  const [tokens, setTokens] = useState([]);
  const siteInfo = useSelector((state) => state.siteInfo);

  const fetchData = async () => {
    try {
      const res = await API.get(`/api/token/`, {
        params: {
          page: 1,
          size: 10,
          keyword: '',
          order: '-id'
        }
      });
      const { success, message, data } = res.data;
      if (success) {
        setTokens(data.data);
      } else {
        showError(message);
      }
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);
  let serverAddress = '';
  if (siteInfo?.server_address) {
    serverAddress = siteInfo.server_address;
  } else {
    serverAddress = window.location.host;
  }
  return (
    <>
      {/* logo & toggler button */}
      <Box
        sx={{
          width: 228,
          display: 'flex',
          [theme.breakpoints.down('md')]: {
            width: 'auto'
          }
        }}
      >
        <Box component="span" sx={{ display: { xs: 'none', md: 'block' }, flexGrow: 1 }}>
          <LogoSection />
        </Box>
        <ButtonBase sx={{ borderRadius: '12px', overflow: 'hidden' }}>
          <Avatar
            variant="rounded"
            sx={{
              ...theme.typography.commonAvatar,
              ...theme.typography.mediumAvatar,
              transition: 'all .2s ease-in-out',
              background: theme.palette.secondary.light,
              color: theme.palette.secondary.dark,
              '&:hover': {
                background: theme.palette.secondary.dark,
                color: theme.palette.secondary.light
              }
            }}
            onClick={handleLeftDrawerToggle}
            color="inherit"
          >
            <IconMenu2 stroke={1.5} size="1.3rem" />
          </Avatar>
        </ButtonBase>
      </Box>
      <Box sx={{ flexGrow: 1 }} />
      <Box sx={{ textAlign: 'center' }}>
        <Button
          sx={{ marginRight: '10px', borderRadius: '15px' }}
          variant="contained"
          href={`https://www.hustgpt.com/#/?settings={"key":"sk-${tokens[0] && tokens[0].key}","url":"${serverAddress}"}`}
          color="primary"
        >
          去聊天
        </Button>
      </Box>
      <Box sx={{ textAlign: 'center' }}>
        <Button
          sx={{ marginRight: '10px', borderRadius: '15px' }}
          variant="contained"
          href={`/about`}
          color="primary"
        >
          使用简介
        </Button>
      </Box>

      <ProfileSection />
    </>
  );
};

Header.propTypes = {
  handleLeftDrawerToggle: PropTypes.func
};

export default Header;

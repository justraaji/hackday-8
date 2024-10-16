import { useState } from 'react'
import { Mosaic } from 'react-loading-indicators'
import { useSpring, animated } from 'react-spring';

const userFields = [
  "Username",
  "Firstname",
  "Lastname",
  "Email",
  "Phone",
]

const companyFields = [
  "Name",
  "Industry",
  "Address",
  "City",
  "State",
  "Zipcode",
]

const InputField = ({ labelName, handler, disabled }) => {
  return (
    <div className='w-full mb-4 flex justify-between items-start flex-col'>
      <label className='text-sm mb-1 font-bold'>
       {labelName}
      </label>
      <input 
        type="text"
        disabled={disabled}
        name={labelName.toLowerCase()}
        className='border-[1px] border-gray-300 rounded-[4px] p-3 w-[400px]'
        onChange={handler}
      />
    </div>
  )
}

const Title = ({title}) => (
  <p className='font-bold text-2xl mb-8'>{title}</p>
)

const UserFields = ({ userFieldHandler }) => (
  <div>
    <Title title={"User Info"} />
    {
      userFields.map((field, indx) => (
        <InputField key={indx} labelName={field} handler={userFieldHandler} />
      ))
    }
  </div>
)

const CompanyFields = ({ companyFieldHandler }) => (
  <div>
    <Title title={"Company Info"} />
    {
      companyFields.map((field, indx) => (
        <InputField key={indx} labelName={field} handler={companyFieldHandler} />
      ))
    }
  </div>
)

const LoadingIcon = () => (
  <div className='flex flex-col justify-center items-center h-screen'>
    <Mosaic color={["#33CCCC", "#33CC36", "#B8CC33", "#FCCA00"]} size="large" text="" textColor="" />
    <p className='my-4 text-xl font-bold'>LOADING RESULTS...</p>
  </div>
)

export const Fields = ({ isLoading, setIsLoading, handleDataSubmission }) => {
  const [userfields, setUserFields] = useState({})
  const [companyfields, setCompanyFields] = useState({})
  const [showUserForm, setShowUserForm] = useState(true)

  const toggleComponents = () => {
    setShowUserForm((prev) => !prev);
  }

  const slideInA = useSpring({
    opacity: showUserForm ? 1 : 0,
    transform: showUserForm ? 'translateX(0%)' : 'translateX(-100%)',
    config: { tension: 300, friction: 30 },
  });

  const slideInB = useSpring({
    opacity: showUserForm ? 0 : 1,
    transform: showUserForm ? 'translateX(100%)' : 'translateX(0%)',
    config: { tension: 300, friction: 30 },
  });
  

  const handleUserField = (e) => {
    setUserFields((prev) => ({
      ...prev,
      [e.target.name]: e.target.value
    }))
  }
  
  const handleCompanyField = (e) => {
    let prefix = "company_";
    let key = prefix + e.target.name

    setCompanyFields((prev) => ({
      ...prev,
      [key]: e.target.value
    }))
  }

  const handleBtnClick = () => {
    if (showUserForm) {
        toggleComponents();
    }

    if (!showUserForm) {
      setIsLoading(true)
      let DTO = {
        ...userfields,
        ...companyfields
      }
      handleDataSubmission(DTO)
    }
  }

  return (
    <>
      <div className='flex justify-center items-center'>
          {
            showUserForm 
            ? <animated.div style={slideInA}>
                <UserFields userFieldHandler={handleUserField} />
              </animated.div> 
            : isLoading 
            ? <LoadingIcon />
            : <animated.div style={slideInB}>
                <CompanyFields companyFieldHandler={handleCompanyField} />
              </animated.div>
          }
      </div>
      {
        isLoading 
        ? null
        : <div className='flex justify-center py-3'>
            <button 
              onClick={handleBtnClick}
              className='bg-black text-white py-4 rounded-sm w-[300px] hover:bg-blue-600 transition-all'
            >
            {
              showUserForm 
              ? "Continue"
              : "SUBMIT"
            }
            </button>
          </div>
      }
    </>
  )
}

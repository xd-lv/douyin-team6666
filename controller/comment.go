package controller

import (
	"github.com/gin-gonic/gin"
	"main/constants"
	"main/pack"
	"main/service"
	"main/utils/jwtUtil"
	"net/http"
	"strconv"
)

type CommentActionRequest struct {
	UserId      int64  `json:"user_id,omitempty"`
	VideoId     int64  `json:"video_id,omitempty"`
	ActionType  int32  `json:"action_type,omitempty"`
	Token       string `json:"token,omitempty"`
	CommentText string `json:"comment_text"`
	CommentId   int64  `json:"comment_id"`
}

type CommentActionResponse struct {
	pack.Response
	Comment pack.Comment `json:"comment"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	commentAction := CommentActionRequest{}
	var err error
	commentAction.CommentId, err = strconv.ParseInt(c.Query("comment_id"), 10, 64)
	commentAction.CommentText = c.Query("comment_text")
	tempActionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	commentAction.ActionType = int32(tempActionType)
	commentAction.VideoId, err = strconv.ParseInt(c.Query("video_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: err.Error()},
			Comment:  pack.Comment{},
		})
		return
	}

	claim, err := jwtUtil.AuthMiddleware.GetClaimsFromJWT(c)
	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1, StatusMsg: "User doesn't exist",
		})
		return
	}
	userId, err := strconv.ParseInt(claim[constants.IdentityKey].(string), 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, pack.Response{
			StatusCode: 1, StatusMsg: "User doesn't exist~~~~" + strconv.FormatInt(userId, 10) + " " + strconv.FormatInt(commentAction.UserId, 10),
		})
		return
	}
	commentAction.UserId = userId

	commentAction.Token = c.Query("token")
	if commentAction.ActionType == constants.ActionPublishComment {
		comment, err := service.CommentService.CreateComment(c, commentAction.VideoId, commentAction.UserId, commentAction.CommentText)

		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: pack.Response{StatusCode: 2, StatusMsg: err.Error()},
				Comment:  pack.Comment{},
			})
			return
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: pack.Response{StatusCode: 0},
			Comment:  comment,
		})
		return
	} else if commentAction.ActionType == constants.ActionDeleteComment {
		err := service.CommentService.DeleteComment(c, commentAction.VideoId, commentAction.CommentId)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: pack.Response{StatusCode: 2, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: pack.Response{StatusCode: 0},
		})

	} else {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: constants.ErrInvalidParams.Error()},
		})
		return
	}

}

type CommentListResponse struct {
	pack.Response
	CommentList []pack.Comment `json:"comment_list"`
}

type CommentListRequest struct {
	VideoId int64  `json:"video_id,omitempty"`
	Token   string `json:"token,omitempty"`
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	commentListRequest := CommentListRequest{}
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	commentListRequest.VideoId = videoId
	commentListRequest.Token = c.Query("token")
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	claim, err := jwtUtil.AuthMiddleware.GetClaimsFromJWT(c)

	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: "User doesn't exist" + err.Error()},
		})
		return
	}

	_, err = strconv.ParseInt(claim[constants.IdentityKey].(string), 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: "User doesn't exist" + err.Error()},
		})
		return
	}
	comment, err := service.CommentService.ListComment(c, commentListRequest.VideoId)

	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: pack.Response{StatusCode: 1, StatusMsg: "video doesn't exist" + err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    pack.Response{StatusCode: 0},
		CommentList: comment,
	})

}
